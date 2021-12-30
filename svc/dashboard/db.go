package main

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"

	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/password"
	"lab.weave.nl/nid/nid-core/svc/auth/models"
)

// DashboardDB database struct used by dashboard grpc services
type DashboardDB struct {
	db     *gorm.DB
	UserDB *models.UserDB
}

func initDB(conf *DashBoardConfig) *DashboardDB {
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

	dashBoardDB := &DashboardDB{
		db:     db,
		UserDB: models.NewUserDB(db),
	}
	dashBoardDB.migrate(conf)

	return dashBoardDB
}

func (d *DashboardDB) migrate(conf *DashBoardConfig) {
	migrationList := []*gormigrate.Migration{
		{
			ID: "2019-11-26-ADD-DEFAULT-USER",
			Migrate: func(tx *gorm.DB) error {
				m := password.NewDefaultManager()
				password, err := m.GenerateHash(conf.DefaultUserPass)
				if err != nil {
					return errors.Wrap(err, "unable to generate hash for default user")
				}
				err = tx.Create(&models.User{
					Email:    conf.DefaultUser,
					Password: password,
				}).Error
				if err != nil {
					return errors.Wrap(err, "unable to create default user")
				}

				return nil
			},
		},
		{
			ID: "2019-11-09-ADD-PILOT-USER",
			Migrate: func(tx *gorm.DB) error {
				m := password.NewDefaultManager()
				password, err := m.GenerateHash(conf.PilotUserPass)
				if err != nil {
					return errors.Wrap(err, "unable to generate hash for pilot user")
				}
				err = tx.Create(&models.User{
					Email:    conf.PilotUser,
					Password: password,
				}).Error
				if err != nil {
					return errors.Wrap(err, "unable to create pilot user")
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
