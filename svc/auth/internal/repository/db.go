// Package repository contains the code about storages (postgres)
package repository

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/config"
	"gopkg.in/gormigrate.v1"

	"github.com/nID-sourcecode/nid-core/pkg/password"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/database/v2"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
)

const (
	NoPreload                        models.PreloadOption = 0
	PreloadRequiredAndOptionalScopes models.PreloadOption = 1
	PreloadAll                       models.PreloadOption = 2
)

// AuthDB database struct used by auth grpc services
type AuthDB struct {
	DB               *gorm.DB
	AccessModelDB    *models.AccessModelDB
	AudienceDB       *models.AudienceDB
	ClientDB         *models.ClientDB
	RedirectTargetDB *models.RedirectTargetDB
	SessionDB        *models.SessionDB
	UserDB           *models.UserDB
	RefreshTokenDB   *models.RefreshTokenDB
	ScopeDB          *models.ScopeDB
}

func InitDB(conf *config.AuthConfig) *AuthDB {
	db := database.MustConnectCustomWithCustomLogger(&database.DBConfig{
		Host:           conf.PGHost,
		User:           conf.PGUser,
		Pass:           conf.PGPass,
		Port:           conf.PGPort,
		RetryOnFailure: true,
		TestMode:       database.TestModeOff,
		DBName:         "auth",
		LogMode:        false,
		AutoMigrate:    true,
		Extensions:     []string{"uuid-ossp"},
	}, models.GetModels(), log.GetLogger())

	authDB := NewAuthDB(db)
	authDB.migrate(conf)

	return authDB
}

func NewAuthDB(db *gorm.DB) *AuthDB {
	return &AuthDB{
		DB:               db,
		AccessModelDB:    models.NewAccessModelDB(db),
		AudienceDB:       models.NewAudienceDB(db),
		ClientDB:         models.NewClientDB(db),
		RedirectTargetDB: models.NewRedirectTargetDB(db),
		SessionDB:        models.NewSessionDB(db),
		UserDB:           models.NewUserDB(db),
		RefreshTokenDB:   models.NewRefreshTokenDB(db),
		ScopeDB:          models.NewScopeDB(db),
	}
}

//nolint:gocognit,funlen,gocyclo
func (d AuthDB) migrate(conf *config.AuthConfig) {
	migrationList := []*gormigrate.Migration{
		{
			ID: "2020-10-02-SEED-AUTH",
			Migrate: func(tx *gorm.DB) error {
				client := d.ClientDB.DefaultModel()
				err := tx.Create(&client).Error
				if err != nil {
					return err
				}

				redirectTarget := d.RedirectTargetDB.DefaultModel(client.ID)
				err = tx.Create(&redirectTarget).Error
				if err != nil {
					return err
				}

				audience := d.AudienceDB.DefaultModelPrimary(conf.Namespace)
				err = tx.Create(&audience).Error
				if err != nil {
					return err
				}

				accessModel := d.AccessModelDB.DefaultModelUserFirstAndLastName(conf.Namespace)
				err = tx.Create(&accessModel).Error
				if err != nil {
					return err
				}
				accessModelopt := d.AccessModelDB.DefaultModelOptionalBankAccounts(conf.Namespace)
				err = tx.Create(&accessModelopt).Error
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			ID: "2020-10-07-SEED-ADD-AUDIENCE-AND-TWO-NEW-ACCESS-MODELS",
			Migrate: func(tx *gorm.DB) error {
				audience := d.AudienceDB.DefaultModelSecondary(conf.Namespace)
				err := tx.Create(&audience).Error
				if err != nil {
					return err
				}
				accessModel := d.AccessModelDB.DefaultModelUserAddresses(conf.Namespace)
				err = tx.Create(&accessModel).Error
				if err != nil {
					return err
				}
				accessModelopt := d.AccessModelDB.DefaultModelOptionalUserAddressContactDetails(conf.Namespace)
				err = tx.Create(&accessModelopt).Error
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			ID: "2020-10-22-CHANGE-SECONDARY-AUDIENCE-URL",
			Migrate: func(tx *gorm.DB) error {
				a := d.AudienceDB.DefaultModelSecondary(conf.Namespace)
				err := tx.Model(&a).Where("id = ?", a.ID).Update("audience", a.Audience).Error
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			ID: "2020-11-11-SEED-PILOT-CLIENT",
			Migrate: func(tx *gorm.DB) error {
				client := models.Client{
					ID:   uuid.Must(uuid.FromString("539e94b9-c5bb-4025-b325-8efa2be2d75d")),
					Name: "Pilot Client",
				}
				err := tx.Create(&client).Error
				if err != nil {
					return err
				}

				redirectTarget := d.RedirectTargetDB.DefaultModel(client.ID)
				err = tx.Create(&redirectTarget).Error
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			ID: "2020-11-12-MIGRATE-GQL-ACCESS-MODELS-TO-NEW-TABLE",
			Migrate: func(tx *gorm.DB) error {
				oldStyleAccessModels := make([]*models.AccessModel, 0)
				err := tx.Find(&oldStyleAccessModels, "type IS NULL").Error
				log.Infof("models: %d", len(oldStyleAccessModels))
				if err != nil {
					return errors.Wrap(err, "finding access models with nil types")
				}

				for _, accessModel := range oldStyleAccessModels {
					accessModel.Type = models.AccessModelTypeGQL
					accessModel.GqlAccessModel = &models.GqlAccessModel{
						JSONModel: accessModel.JSONModel,
						Path:      "/gql",
					}
					err = tx.Save(accessModel).Error
					if err != nil {
						return errors.Wrap(err, "updating access model with nil type to GQL")
					}
				}

				return nil
			},
		},
		{
			ID: "2020-11-13-SEED-INFORMATION-SERVICE-AUTH",
			Migrate: func(tx *gorm.DB) error {
				audience := d.AudienceDB.DefaultModelInformationService(conf.Namespace)
				err := tx.Create(&audience).Error
				if err != nil {
					return errors.Wrap(err, "inserting information service audience")
				}

				accessModelContactable := d.AccessModelDB.DefaultModelContactable(conf.Namespace)
				err = tx.Create(&accessModelContactable).Error
				if err != nil {
					return errors.Wrap(err, "inserting information service audience")
				}

				accessModelHasAddress := d.AccessModelDB.DefaultModelHasAddress(conf.Namespace)
				err = tx.Create(&accessModelHasAddress).Error
				if err != nil {
					return errors.Wrap(err, "inserting information service audience")
				}

				accessModelPositiveBalance := d.AccessModelDB.DefaultModelHasPositiveBankAccountBalance(conf.Namespace)
				err = tx.Create(&accessModelPositiveBalance).Error
				if err != nil {
					return errors.Wrap(err, "inserting information service audience")
				}

				return nil
			},
		},
		{
			ID: "2020-11-13-SET-DEFAULT-CLIENT-PASSWORD",
			Migrate: func(tx *gorm.DB) error {
				m := password.NewDefaultManager()
				testingPasswordHash, err := m.GenerateHash(conf.TestingClientPassword)
				if err != nil {
					return errors.Wrap(err, "generating testing client password hash")
				}

				client := d.ClientDB.DefaultModel()
				client.Password = testingPasswordHash
				err = tx.Save(&client).Error
				if err != nil {
					return errors.Wrap(err, "updating testing client password")
				}

				pilotClientHash, err := m.GenerateHash(conf.PilotClientPassword)
				if err != nil {
					return errors.Wrap(err, "generating pilot client password hash")
				}

				pilotClient := &models.Client{}
				tx.Find(pilotClient, "id = ?", "539e94b9-c5bb-4025-b325-8efa2be2d75d")
				pilotClient.Password = pilotClientHash
				err = tx.Save(&pilotClient).Error
				if err != nil {
					return errors.Wrap(err, "updating pilot client password")
				}

				return nil
			},
		},
		{
			ID: "2020-05-20-UPDATE-DEFAULT-CLIENT-METADATA",
			Migrate: func(tx *gorm.DB) error {
				defaultClient := d.ClientDB.DefaultModel()
				existingClient := &models.Client{
					ID: defaultClient.ID,
				}
				err := tx.Find(existingClient).Error
				if err != nil {
					return errors.Wrap(err, "getting existing default client")
				}

				existingClient.Metadata = defaultClient.Metadata
				err = tx.Save(existingClient).Error
				return errors.Wrap(err, "updating client with metadata")
			},
		},
		{
			ID: "2020-07-23-ADD-BSN-ACCESS-MODELS",
			Migrate: func(tx *gorm.DB) error {
				accessModel := d.AccessModelDB.DefaultModelUserFirstAndLastNameByBSN(conf.Namespace)
				err := tx.Create(&accessModel).Error
				if err != nil {
					return errors.Wrap(err, "creating name (by bsn) access model")
				}
				accessModelopt := d.AccessModelDB.DefaultModelOptionalBankAccountsByBSN(conf.Namespace)
				err = tx.Create(&accessModelopt).Error

				return errors.Wrap(err, "creating bank account (by bsn) access model")
			},
		},
	}

	m := gormigrate.New(d.DB, gormigrate.DefaultOptions, migrationList)
	var err error
	if len(migrationList) > 0 {
		err = m.Migrate()
	}
	if err != nil {
		log.WithError(err).Fatal("unable to migrate db")
	}
	log.Info("Migration ran successfully")
}
