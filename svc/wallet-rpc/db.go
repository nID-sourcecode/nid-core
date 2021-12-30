package main

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"

	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/password"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/models"
)

// WalletDB database struct used by wallet grpc services
type WalletDB struct {
	db             *gorm.DB
	UserDB         *models.UserDB
	DeviceDB       *models.DeviceDB
	ConsentDB      *models.ConsentDB
	ClientDB       *models.ClientDB
	EmailAddressDB *models.EmailAddressDB
	PhoneNumberDB  *models.PhoneNumberDB
}

func initDB(conf *WalletConfig, testMode bool, passManager password.IManager) *WalletDB {
	var db *gorm.DB
	if testMode {
		db = database.MustConnectTest("wallet", models.GetModels())
	} else {
		db = database.MustConnectCustomWithCustomLogger(&database.DBConfig{
			Host:           conf.PGHost,
			User:           conf.PGUser,
			Pass:           conf.PGPass,
			Port:           conf.PGPort,
			RetryOnFailure: true,
			TestMode:       database.TestModeOff,
			DBName:         "wallet",
			LogMode:        false,
			AutoMigrate:    true,
			Extensions:     []string{"uuid-ossp"},
		}, models.GetModels(), log.GetLogger())
	}

	dashBoardDB := &WalletDB{
		db:             db,
		UserDB:         models.NewUserDB(db),
		DeviceDB:       models.NewDeviceDB(db),
		ConsentDB:      models.NewConsentDB(db),
		ClientDB:       models.NewClientDB(db),
		EmailAddressDB: models.NewEmailAddressDB(db),
		PhoneNumberDB:  models.NewPhoneNumberDB(db),
	}
	if !testMode {
		dashBoardDB.migrate(conf, passManager)
	}

	return dashBoardDB
}

type defaultUser struct {
	Bsn       string           `json:"bsn"`
	Password  string           `json:"password"`
	Pseudonym string           `json:"pseudonym"`
	Devices   []*defaultDevice `json:"devices"`
}

type defaultDevice struct {
	Code   string `json:"code"`
	Secret string `json:"secret"`
}

func (d *WalletDB) migrate(conf *WalletConfig, passManager password.IManager) {
	migrationList := []*gormigrate.Migration{
		{
			ID: "2020-11-24-ADD-DEFAULT-USERS",
			Migrate: func(tx *gorm.DB) error {
				defaultUsers := make([]*defaultUser, 0)
				err := json.Unmarshal([]byte(conf.DefaultUsers), &defaultUsers)
				if err != nil {
					return errors.Wrap(err, "parsing default users JSON")
				}

				for _, defaultUser := range defaultUsers {
					pwHash, err := passManager.GenerateHash(defaultUser.Password)
					if err != nil {
						return errors.Wrapf(err, "generating hash for defaultUser %d", defaultUser.Bsn)
					}
					newUser := models.User{
						Bsn:       defaultUser.Bsn,
						Password:  pwHash,
						Pseudonym: defaultUser.Pseudonym,
					}
					err = tx.Create(&newUser).Error
					if err != nil {
						return errors.Wrap(err, "inserting default user")
					}

					if defaultUser.Devices != nil {
						for _, defaultDevice := range defaultUser.Devices {
							secretHash, err := passManager.GenerateHash(defaultDevice.Secret)
							if err != nil {
								return errors.Wrapf(err, "generating hash for defaultDevice with code %d", defaultDevice.Code)
							}
							err = tx.Create(&models.Device{
								UserID: newUser.ID,
								Secret: secretHash,
								Code:   defaultDevice.Code,
							}).Error
							if err != nil {
								return errors.Wrapf(err, `inserting default device "%s"`, defaultDevice.Code)
							}
						}
					}
				}

				return nil
			},
		},
	}
	log.Infof("Migration started")
	m := gormigrate.New(d.db, gormigrate.DefaultOptions, migrationList)
	var err error
	if len(migrationList) > 0 {
		err = m.Migrate()
	}
	if err != nil {
		log.WithError(err).Fatal("unable to migrate db")
	}
	log.Info("Migration ran successfully")
}
