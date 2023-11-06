package models

import (
	"context"
	"testing"
	"time"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/authtoken"
	"github.com/nID-sourcecode/nid-core/pkg/password"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/database/v2"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
)

const databaseName = "auth"

type SessionAddOnTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	db              *gorm.DB
	tx              *gorm.DB
	SessionDB       *SessionDB
	ClientDB        *ClientDB
	session         *Session
	passwordManager password.IManager
	audience        *Audience
	client          *Client
}

func (s *SessionAddOnTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()

	s.passwordManager = password.NewDefaultManager()

	s.tx = s.db.BeginTx(context.TODO(), nil)
	s.Require().NoError(s.tx.AutoMigrate(GetModels()...).Error)

	// s.Require().Len(AddForeignKeys(s.tx), 0) FIXME reintroduce foreign keys when generator bug is fixed

	s.SessionDB = NewSessionDB(s.tx)
	s.ClientDB = NewClientDB(s.tx)
	session, err := s.createDummySession()
	s.Require().NoError(err)
	s.session = session
}

func (s *SessionAddOnTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *SessionAddOnTestSuite) TearDownSuite() {
	s.NoError(s.db.Close())
	log.Info("Closing db connection")
}

func (s *SessionAddOnTestSuite) createDummySession() (*Session, error) {
	client1Password, err := s.passwordManager.GenerateHash("test^123")
	s.Require().NoError(err)
	client2Password, err := s.passwordManager.GenerateHash("456#%test")
	s.Require().NoError(err)
	client := &Client{
		Color:    "blue",
		Name:     "testclient",
		Password: client1Password,
	}
	err = s.tx.Create(client).Error
	if err != nil {
		return nil, err
	}

	s.client = client

	client2 := &Client{
		Color:    "red",
		Name:     "testclient2",
		Password: client2Password,
	}
	err = s.tx.Create(client2).Error
	if err != nil {
		return nil, err
	}

	redirectTarget := &RedirectTarget{
		ClientID:       client.ID,
		RedirectTarget: "https://weave.nl/code",
	}
	err = s.tx.Create(redirectTarget).Error
	if err != nil {
		return nil, err
	}

	redirectTarget2 := &RedirectTarget{
		ClientID:       client2.ID,
		RedirectTarget: "https://weave2.nl/code",
		UpdatedAt:      time.Time{},
	}
	err = s.tx.Create(redirectTarget2).Error
	if err != nil {
		return nil, err
	}

	s.audience = &Audience{
		Audience:  "https://test.com/gql",
		Namespace: "alice",
	}
	err = s.tx.Create(s.audience).Error
	if err != nil {
		return nil, err
	}

	audience2 := &Audience{
		Audience:  "https://test2.com/gql",
		Namespace: "bob",
	}
	err = s.tx.Create(audience2).Error
	if err != nil {
		return nil, err
	}

	accessModels := []*AccessModel{
		{
			AudienceID: s.audience.ID,
			Hash:       "abc",
			Name:       "test:stuff",
			Type:       AccessModelTypeGQL,
			GqlAccessModel: &GqlAccessModel{
				Path:      "/gql",
				JSONModel: `{"r":"somestuff"}`,
			},
		},
	}

	session := &Session{
		ID:                   client.ID,
		AudienceID:           s.audience.ID,
		ClientID:             client.ID,
		RedirectTargetID:     redirectTarget.ID,
		AcceptedAccessModels: accessModels,
	}

	err = s.tx.Create(session).Error
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionAddOnTestSuite) TestCreateSession() {
	// Session is created in SetupTest since it is also used by other tests
	var foundSession Session
	err := s.tx.Where("id = ?", s.session.ID).Find(&foundSession).Error
	s.Require().NoError(err)
	s.Equal(s.session.ClientID, foundSession.ClientID)
	s.Equal(s.session.RedirectTargetID, foundSession.RedirectTargetID)
}

func (s *SessionAddOnTestSuite) TestUpdateAcceptedAccessModels() {
	var acceptedAccessModel []*AccessModel
	err := s.tx.Model(&s.session).Association("AcceptedAccessModels").Find(&acceptedAccessModel).Error
	s.Require().NoError(err)
	s.Len(acceptedAccessModel, 1)

	s.session.AcceptedAccessModels = append(s.session.AcceptedAccessModels, &AccessModel{
		AudienceID: s.audience.ID,
		Hash:       "ghi",
		Name:       "test:stuff2",
		Type:       AccessModelTypeGQL,
		GqlAccessModel: &GqlAccessModel{
			Path:      "/graphql",
			JSONModel: `{"r":"somemorestuff"}`,
		},
	})
	err = s.SessionDB.UpdateAcceptedAccessModels(s.session, s.session.AcceptedAccessModels)
	s.NoError(err)

	var updatedAcceptedAccessModel []*AccessModel
	err = s.tx.Model(&s.session).Association("AcceptedAccessModels").Find(&updatedAcceptedAccessModel).Error
	s.NoError(err)
	s.Len(updatedAcceptedAccessModel, len(s.session.AcceptedAccessModels))
}

func (s *SessionAddOnTestSuite) TestUpdateSessionState() {
	s.Require().NotEqual(s.session.State, SessionStateAccepted)
	err := s.SessionDB.UpdateSessionState(s.session, SessionStateAccepted)
	s.NoError(err)
	s.Equal(s.session.State, SessionStateAccepted)
}

func (s *SessionAddOnTestSuite) TestUpdateSessionStateUpdateOnlyState() {
	s.Require().NotEqual(s.session.State, SessionStateAccepted)

	// Setup
	originalSession := s.session
	testClientColor := s.client.Color

	// Update test client
	s.client.Color = "purple"

	// Update the session
	s.session = &Session{
		Client: s.client,
	}

	// UpdateSessionState
	err := s.SessionDB.UpdateSessionState(s.session, SessionStateAccepted)
	s.Require().NoError(err)
	s.Require().Equal(s.session.State, SessionStateAccepted)

	// Check if client was updated
	databaseTestClient := &Client{}
	s.ClientDB.Db.Where("id = ?", s.client.ID).First(databaseTestClient)

	s.NotEqual(databaseTestClient.ID, uuid.Nil, "testclient-session-addon not created in database")
	s.Equal(testClientColor, databaseTestClient.Color, "Client color should not have changed")

	// Cleanup
	s.session = originalSession
}

func (s *SessionAddOnTestSuite) TestUpdateSessionSubject() {
	subject := "sadasdasjkdhaiouysdg867ig672315471r23t7"
	err := s.SessionDB.UpdateSessionSubject(s.session, subject)
	s.Require().NoError(err)
	s.Require().Equal(s.session.Subject, subject)
}

func (s *SessionAddOnTestSuite) TestUpdateAuthorizationCode() {
	token, err := authtoken.NewToken(32)
	s.Require().NoError(err)
	hash, err := authtoken.Hash(token)
	s.Require().NoError(err)
	err = s.SessionDB.UpdateSessionAuthorizationCode(s.session, hash)
	s.Require().NoError(err)
	s.Require().Equal(s.session.AuthorizationCode, &hash)
}

func TestSessionAddOnTestSuite(t *testing.T) {
	suite.Run(t, &SessionAddOnTestSuite{
		// Intentionally do not supply models to automigrate, this should be done inside the transaction
		db: database.MustConnectTest(databaseName, nil),
	})
}
