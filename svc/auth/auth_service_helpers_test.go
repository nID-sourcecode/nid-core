package main

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/svc/auth/models"
	pb "lab.weave.nl/nid/nid-core/svc/auth/proto"
)

type HelpersTestSuite struct {
	AuthServiceBaseTestSuite
}

func (s *HelpersTestSuite) SetupTest() {
	s.AuthServiceBaseTestSuite.SetupTest()
}

func (s *HelpersTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *HelpersTestSuite) SetupSuite() {
	s.AuthServiceBaseTestSuite.SetupSuite()
}

func (s *HelpersTestSuite) TearDownSuite() {
	s.Require().NoError(s.db.Close())
}

func (s *HelpersTestSuite) TestModelToResponse() {
	session := &models.Session{
		ID: uuid.FromStringOrNil("4488a760-2e90-47f8-a6ac-8e5aaa4985fa"),
		AcceptedAccessModels: []*models.AccessModel{
			s.accessModelGql1,
		},
		OptionalAccessModels: []*models.AccessModel{
			s.accessModelGql1,
		},
		RequiredAccessModels: []*models.AccessModel{
			s.accessModelGql2,
			s.accessModelRestPOST,
		},
		Audience: s.audience,
		Client:   s.client,
		State:    models.SessionStateCodeGranted,
	}
	response := sessionToResponse(session)
	s.Equal(response.Id, uuid.FromStringOrNil("4488a760-2e90-47f8-a6ac-8e5aaa4985fa").String())
	s.Equal(response.Audience.Id, s.audience.ID.String())
	s.Equal(response.Audience.Audience, s.audience.Audience)
	s.Equal(response.Audience.Namespace, s.audience.Namespace)
	s.Equal(response.Client.Id, s.client.ID.String())
	s.Equal(response.Client.Name, s.client.Name)
	s.Equal(response.Client.Color, s.client.Color)
	s.Equal(response.Client.Logo, s.client.Logo)
	s.Equal(response.Client.Icon, s.client.Icon)
	s.Len(response.AcceptedAccessModels, len([]*models.AccessModel{s.accessModelGql1}))
	s.Len(response.OptionalAccessModels, len([]*models.AccessModel{s.accessModelGql1}))
	s.Len(response.RequiredAccessModels, len([]*models.AccessModel{s.accessModelGql2, s.accessModelRestPOST}))
	if s.NotNil(response.AcceptedAccessModels[0]) {
		s.Equal(response.AcceptedAccessModels[0].Id, s.accessModelGql1.ID.String())
		s.Equal(response.AcceptedAccessModels[0].Name, s.accessModelGql1.Name)
		s.Equal(response.AcceptedAccessModels[0].Description, s.accessModelGql1.Description)
		s.Equal(response.AcceptedAccessModels[0].Hash, s.accessModelGql1.Hash)
	}
	s.Equal(response.State, pb.SessionState_CODE_GRANTED)
}

func (s *HelpersTestSuite) TestGetState() {
	tests := []struct {
		Name          string
		GivenState    models.SessionState
		ExpectedState pb.SessionState
		ErrorExpected bool
	}{
		{
			Name:          "models.SessionStateUnclaimed to pb.SessionState_Unclaimed",
			GivenState:    models.SessionStateUnclaimed,
			ExpectedState: pb.SessionState_UNCLAIMED,
			ErrorExpected: false,
		},
		{
			Name:          "models.SessionStateClaimed to pb.SessionState_Claimed",
			GivenState:    models.SessionStateClaimed,
			ExpectedState: pb.SessionState_CLAIMED,
			ErrorExpected: false,
		},
		{
			Name:          "models.SessionStateAccepted to pb.SessionState_Accepted",
			GivenState:    models.SessionStateAccepted,
			ExpectedState: pb.SessionState_ACCEPTED,
			ErrorExpected: false,
		},
		{
			Name:          "models.SessionStateRejected to pb.SessionState_Rejected",
			GivenState:    models.SessionStateRejected,
			ExpectedState: pb.SessionState_REJECTED,
			ErrorExpected: false,
		},
		{
			Name:          "models.SessionStateCodeGranted to pb.SessionState_CodeGranted",
			GivenState:    models.SessionStateCodeGranted,
			ExpectedState: pb.SessionState_CODE_GRANTED,
			ErrorExpected: false,
		},
		{
			Name:          "models.SessionStateTokenGranted to pb.SessionState_TokenGranted",
			GivenState:    models.SessionStateTokenGranted,
			ExpectedState: pb.SessionState_TOKEN_GRANTED,
			ErrorExpected: false,
		},
		{
			Name:          "models.SessionStateUnclaimed to pb.SessionState_TokenGranted",
			GivenState:    models.SessionStateUnclaimed,
			ExpectedState: pb.SessionState_TOKEN_GRANTED,
			ErrorExpected: true,
		},
	}
	for _, test := range tests {
		s.Run(test.Name, func() {
			result := getState(test.GivenState)
			if test.ErrorExpected {
				s.NotEqual(test.ExpectedState, result)
			} else {
				s.Equal(test.ExpectedState, result)
			}
		})
	}
}

func (s *HelpersTestSuite) TestFindMultipleAccessModels() {
	optionalModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("4c6b94c0-e6e9-478a-b19e-2e549be72fae"),
		},
		{
			ID: uuid.FromStringOrNil("dd6a8ce4-f8be-4070-8a7c-a30d12349fe1"),
		},
		{
			ID: uuid.FromStringOrNil("1a004d41-59b1-4489-8172-42b28630b274"),
		},
	}
	requiredModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("c09d687b-ba94-412e-af17-84cfc76f3f26"),
		},
		{
			ID: uuid.FromStringOrNil("72530df1-9f00-4cbe-86c3-d0d31ed2278a"),
		},
	}
	suppliedModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("dd6a8ce4-f8be-4070-8a7c-a30d12349fe1"),
		},
		{
			ID: uuid.FromStringOrNil("1a004d41-59b1-4489-8172-42b28630b274"),
		},
	}
	s.NoError(accessModelExists(suppliedModels, requiredModels, optionalModels))
}

func (s *HelpersTestSuite) TestFindSingleAccessModel() {
	optionalModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("4c6b94c0-e6e9-478a-b19e-2e549be72fae"),
		},
		{
			ID: uuid.FromStringOrNil("dd6a8ce4-f8be-4070-8a7c-a30d12349fe1"),
		},
		{
			ID: uuid.FromStringOrNil("1a004d41-59b1-4489-8172-42b28630b274"),
		},
	}
	requiredModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("c09d687b-ba94-412e-af17-84cfc76f3f26"),
		},
		{
			ID: uuid.FromStringOrNil("72530df1-9f00-4cbe-86c3-d0d31ed2278a"),
		},
	}
	suppliedModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("dd6a8ce4-f8be-4070-8a7c-a30d12349fe1"),
		},
	}
	s.NoError(accessModelExists(suppliedModels, requiredModels, optionalModels))
}

func (s *HelpersTestSuite) TestErrorSupplyOnlyRequiredInPayload() {
	var optionalModels []*models.AccessModel
	requiredModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("c09d687b-ba94-412e-af17-84cfc76f3f26"),
		},
		{
			ID: uuid.FromStringOrNil("72530df1-9f00-4cbe-86c3-d0d31ed2278a"),
		},
	}
	suppliedModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("72530df1-9f00-4cbe-86c3-d0d31ed2278a"),
		},
	}
	s.Error(accessModelExists(suppliedModels, requiredModels, optionalModels))
}

func (s *HelpersTestSuite) TestErrorFindSingleNonExisting() {
	optionalModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("4c6b94c0-e6e9-478a-b19e-2e549be72fae"),
		},
		{
			ID: uuid.FromStringOrNil("dd6a8ce4-f8be-4070-8a7c-a30d12349fe1"),
		},
		{
			ID: uuid.FromStringOrNil("1a004d41-59b1-4489-8172-42b28630b274"),
		},
	}
	requiredModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("c09d687b-ba94-412e-af17-84cfc76f3f26"),
		},
		{
			ID: uuid.FromStringOrNil("72530df1-9f00-4cbe-86c3-d0d31ed2278a"),
		},
	}
	suppliedModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("0f8fc5a5-2c48-4663-972b-34014ae24513"),
		},
	}
	s.Error(accessModelExists(suppliedModels, requiredModels, optionalModels))
}

func (s *HelpersTestSuite) TestErrorFindSingleNonExistingAndSingleExisting() {
	optionalModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("4c6b94c0-e6e9-478a-b19e-2e549be72fae"),
		},
		{
			ID: uuid.FromStringOrNil("dd6a8ce4-f8be-4070-8a7c-a30d12349fe1"),
		},
		{
			ID: uuid.FromStringOrNil("1a004d41-59b1-4489-8172-42b28630b274"),
		},
	}
	requiredModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("c09d687b-ba94-412e-af17-84cfc76f3f26"),
		},
		{
			ID: uuid.FromStringOrNil("72530df1-9f00-4cbe-86c3-d0d31ed2278a"),
		},
	}
	suppliedModels := []*models.AccessModel{
		{
			ID: uuid.FromStringOrNil("72530df1-9f00-4cbe-86c3-d0d31ed2278a"),
		},
		{
			ID: uuid.FromStringOrNil("0f8fc5a5-2c48-4663-972b-34014ae24513"),
		},
	}
	s.Error(accessModelExists(suppliedModels, requiredModels, optionalModels))
}

func TestHelpersTestSuite(t *testing.T) {
	db := database.MustConnectTest("auth", nil)
	fmt.Println("started connection")
	suite.Run(t, &HelpersTestSuite{
		AuthServiceBaseTestSuite: AuthServiceBaseTestSuite{
			db: db,
		},
	})
}
