package design

import (
	"lab.weave.nl/weave/generator"
	. "lab.weave.nl/weave/generator/dsl"
	"lab.weave.nl/weave/generator/types"
)

var _ = Store("auth", generator.Postgres, func() {
	Enum("AccessModelType", func() {
		Value("GQL", 1)
		Value("REST", 2)
	})

	Enum("SessionState", func() {
		Value("unclaimed", 1)
		Value("claimed", 2)
		Value("accepted", 3)
		Value("rejected", 4)
		Value("code_granted", 5)
		Value("token_granted", 6)
	})

	var Client, RedirectTarget, AccessModel, GqlAccessModel, RestAccessModel, Audience, Session, User *generator.RelationalModelDefinition

	_ = Session
	_ = User

	ScopeReadClients := Scope("api:clients:read")

	Client = Model("Client", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		// The client's name.
		Field("Name", types.String())
		// The client's logo.
		Field("Logo", types.String())
		// The client's icon.
		Field("Icon", types.String())
		// The client's color
		Field("Color", types.String())

		Field("Password", types.String()).SetInternalOnly()

		Field("Metadata", types.JSON())

		HasMany(RedirectTarget)

		Read(ScopeReadClients, Relation().None())
	})

	RedirectTarget = Model("RedirectTarget", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("RedirectTarget", types.String())

		BelongsTo(Client)
	})

	AccessModel = Model("AccessModel", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("Name", types.String())
		Field("JSONModel", types.String())
		Field("Hash", types.String())
		Field("Description", types.String())

		BelongsTo(Audience)

		Field("Type", types.Enum("AccessModelType"))
		HasOne(GqlAccessModel)
		HasOne(RestAccessModel)
	})

	GqlAccessModel = Model("GqlAccessModel", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("JSONModel", types.String())
		Field("Path", types.String())

		BelongsTo(AccessModel)
	})

	RestAccessModel = Model("RestAccessModel", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("Path", types.String())
		Field("Query", types.String())
		Field("Body", types.String())
		Field("Method", types.String())

		BelongsTo(AccessModel)
	})

	Audience = Model("Audience", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("Audience", types.String())
		Field("Namespace", types.String())

		HasMany(AccessModel)
	})

	Session = Model("Session", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("State", types.Enum("SessionState"))
		Field("Subject", types.String())
		Field("FinaliseToken", types.String())
		Field("AuthorizationCode", types.String(), func() { Nullable() })
		Field("AuthorizationCodeGrantedAt", types.Timestamp(), func() { Nullable() })
		BelongsTo(Client)
		BelongsTo(Audience)
		BelongsTo(RedirectTarget)
		// ManyToMany("RequiredAccessModels", "AccessModel", "Sessions")
		// ManyToMany("OptionalAccessModels", "AccessModel", "Sessions")
		// ManyToMany("AcceptedAccessModels", "AccessModel", "Sessions")
		ManyToMany(AccessModel, "RequiredAccessModels") //, "AccessModel", "Sessions")
		ManyToMany(AccessModel, "OptionalAccessModels") //, "AccessModel", "Sessions")
		ManyToMany(AccessModel, "AcceptedAccessModels") //, "AccessModel", "Sessions")
	})

	User = Model("User", func() {
		UserDefault(types.UUID())
		Field("ID", types.UUID(), func() { PrimaryKey(); ReadOnly() })
		Field("Email", types.String(), func() { SQLTag("unique") })
		Field("Password", types.String(), func() { WriteOnly(); CustomMutate(types.String()) })
	})
})
