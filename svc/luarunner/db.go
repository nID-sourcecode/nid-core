package main

import (
	"github.com/jinzhu/gorm"

	"lab.weave.nl/nid/nid-core/svc/luarunner/models"
)

// LuaRunnerDB database struct used by LuaRunner service
type LuaRunnerDB struct {
	db             *gorm.DB
	OrganisationDB *models.OrganisationDB
	ScriptDB       *models.ScriptDB
}

// NewLuaRunnerDB returns a new instance of LuaRunnerDB
func NewLuaRunnerDB(db *gorm.DB) *LuaRunnerDB {
	return &LuaRunnerDB{
		db:             db,
		OrganisationDB: models.NewOrganisationDB(db),
		ScriptDB:       models.NewScriptDB(db),
	}
}
