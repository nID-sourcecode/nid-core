package main

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"

	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/svc/info-manager/models"
)

// InfoManagerDB database struct used by info-manager grpc service and info-manager-gql service
type InfoManagerDB struct {
	db             *gorm.DB
	ScriptDB       *models.ScriptDB
	ScriptSourceDB *models.ScriptSourceDB
}

// initDB init database function for info manager service
func initDB(conf *InfoManagerConfig, testMode bool) *InfoManagerDB {
	var db *gorm.DB
	if testMode {
		db = database.MustConnectCustom(&database.DBConfig{
			Host:           conf.PGHost,
			User:           conf.PGUser,
			Pass:           conf.PGPass,
			Port:           conf.PGPort,
			RetryOnFailure: true,
			DBName:         "infomanager",
			TestMode:       database.TestModeNoDropTable,
			AutoMigrate:    true,
			Extensions:     []string{"uuid-ossp"},
		}, models.GetModels())
	} else {
		db = database.MustConnectCustom(&database.DBConfig{
			Host:           conf.PGHost,
			User:           conf.PGUser,
			Pass:           conf.PGPass,
			Port:           conf.PGPort,
			RetryOnFailure: true,
			TestMode:       database.TestModeOff,
			DBName:         "infomanager",
			LogMode:        false,
			AutoMigrate:    true,
			Extensions:     []string{"uuid-ossp"},
		}, models.GetModels())
	}

	infoManagerDB := &InfoManagerDB{
		db:             db,
		ScriptDB:       models.NewScriptDB(db),
		ScriptSourceDB: models.NewScriptSourceDB(db),
	}
	if !testMode {
		infoManagerDB.migrate()
	}

	return infoManagerDB
}

func (d InfoManagerDB) migrate() {
	migrationList := []*gormigrate.Migration{
		{
			ID: "2020-03-09-INITIAL-MIGRATION",
			Migrate: func(tx *gorm.DB) error {
				return nil
			},
		},
	}

	m := gormigrate.New(d.db, gormigrate.DefaultOptions, migrationList)
	var err error
	if len(migrationList) > 0 {
		err = m.Migrate()
	}
	if err != nil {
		log.WithError(err).Fatal("unable to migrate db")
	}

	log.Info("migration ran successfully")
}
