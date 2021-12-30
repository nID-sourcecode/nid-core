package design

import (
	"lab.weave.nl/weave/generator"
	. "lab.weave.nl/weave/generator/dsl"
	"lab.weave.nl/weave/generator/types"
)

var _ = Store("auth", generator.Postgres, func() {
	Enum("PhoneNumberVerificationType", func() {
		Value("SMS", 1)
		Value("TTS", 2)
	})
	var User, EmailAddress, PhoneNumber, Consent, RevokeConsent, Client, Device *generator.RelationalModelDefinition

	_ = User
	_ = EmailAddress
	_ = PhoneNumber
	_ = Consent
	_ = RevokeConsent
	_ = Client
	_ = Device

	User = Model("User", func() {
		UserDefault(types.UUID())
		Field("Pseudonym", types.String(), func() {
			UniqueIndex("pseudonym")
		})
		Field("Bsn", types.String(), func() {
			UniqueIndex("bsn")
		}) // Rainbow table without any security for now

		HasMany(PhoneNumber)
		HasMany(EmailAddress)
		HasMany(Consent)
		HasMany(Device)

		Read(ScopeAccess, Relation().IsMe())
		Create(ScopeAdmin, Relation().None())
		Update(ScopeAdmin, Relation().None())
	})

	Device = Model("Device", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("Code", types.String(), func() {
			UniqueIndex("code")
		})
		Field("Secret", types.String())

		BelongsTo(User)
	})

	EmailAddress = Model("EmailAddress", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })

		Field("EmailAddress", types.String())

		Field("VerificationToken", types.String(), func() { Internal() })
		Field("Verified", types.Boolean(), func() { ReadOnly() })

		BeforeCreate("Hook")
		AfterCreate("Hook")

		BelongsTo(User)

		Read(ScopeAccess, Relation().HasMyUserID())
		Create(ScopeAccess, Relation().HasMyUserID())
		Update(ScopeAccess, Relation().HasMyUserID())
	})

	PhoneNumber = Model("PhoneNumber", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })

		Field("PhoneNumber", types.String())

		Field("VerificationToken", types.String(), func() { Internal() })
		Field("VerificationType", types.Enum("PhoneNumberVerificationType"))
		Field("Verified", types.Boolean(), func() { ReadOnly() })

		BeforeCreate("Hook")
		AfterCreate("Hook")

		BelongsTo(User)

		Read(ScopeAccess, Relation().HasMyUserID())
		Create(ScopeAccess, Relation().HasMyUserID())
		Update(ScopeAccess, Relation().HasMyUserID())
	})

	Consent = Model("Consent", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("Name", types.String())
		Field("Description", types.String())
		Field("Granted", types.Timestamp(), func() { Nullable() })
		Field("Revoked", types.Timestamp(), func() { Nullable() })
		Field("Token", types.JSON())
		Field("AccessToken", types.String())

		BelongsTo(User)
		BelongsTo(Client)

		AfterRead("SetToken")

		Read(ScopeAccess, Relation().HasMyUserID())
	})

	RevokeConsent = Model("RevokeConsent", func() {
		GraphqlOnly()

		Field("ID", types.UUID())
		Field("Revoked", types.Timestamp(), func() { ReadOnly() })

		Create(ScopeAccess, Relation().None())
	})

	Client = Model("Client", func() {
		Field("ID", types.UUID(), func() { PrimaryKey() })
		Field("Name", types.String())
		Field("Logo", types.String())
		Field("Icon", types.String())
		Field("Color", types.String())
		Field("ExtClientId", types.String())

		HasMany(Consent)

		Read(ScopeAccess, Relation().None())
		Create(ScopeAdmin, Relation().None())
		Update(ScopeAdmin, Relation().None())
	})
})
