package design

import (
	"lab.weave.nl/weave/generator"
	. "lab.weave.nl/weave/generator/dsl"
	"lab.weave.nl/weave/generator/types"
)

var _ = Store("infomanager", generator.Postgres, func() {
	Enum("ScriptStatus", func() {
		Value("Draft", 1)
		Value("Active", 2)
		Value("Rejected", 3)
		Value("Archived", 4)
	})

	var Script, ScriptSource *generator.RelationalModelDefinition

	_ = Script
	_ = ScriptSource

	// Script contains some fields with information about a script like name, description and status
	Script = Model("Script", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		// The script's name.
		Field("Name", types.String())
		// The scripts's description.
		Field("Description", types.String())
		// The scripts's status.
		Field("Status", types.Enum("ScriptStatus"))

		HasMany(ScriptSource)

		Read(ScopeOpen, Relation().None())
		Create(ScopeOpen, Relation().None())
		Update(ScopeOpen, Relation().None())
	})

	// Script source contains some fields with information about the script's source like signed_url to source_code file, checksum and version
	ScriptSource = Model("ScriptSource", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		// The source's signed url.
		Field("SignedURL", types.String(), func() { Nullable() })
		// The source's raw script.
		Field("RawScript", types.String(), func() { Nullable() })
		// The source's checksum.
		Field("Checksum", types.String())
		// The source's version.
		Field("Version", types.String())
		// The source's change descriptiob.
		Field("ChangeDescription", types.String())

		AfterRead("GetSignedURL")
		Read(ScopeOpen, Relation().None())
	})
})
