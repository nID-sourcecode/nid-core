package design

import (
	"lab.weave.nl/weave/generator"
	. "lab.weave.nl/weave/generator/dsl"
	"lab.weave.nl/weave/generator/types"
)

var scriptModel *generator.RelationalModelDefinition

var _ = Store("auth", generator.Postgres, func() {
	scriptModel = Model("Script", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("EventType", types.String())
		Field("Script", types.String())

		Create(ScopeOpen, Relation().None())
		Update(ScopeOpen, Relation().None())
		Read(ScopeOpen, Relation().None())
	})

	Model("Organisation", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("Uzovi", types.String())
		Field("RedirectTargetID", types.UUID())
		Field("AudienceID", types.UUID())

		HasMany(scriptModel)
	})
})
