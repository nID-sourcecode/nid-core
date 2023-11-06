package models

import (
	"lab.weave.nl/weave/generator/gen/auth"
)

func GetModels() []interface{} {
	return []interface{}{
		auth.JWT{},
		AccessModel{},
		Audience{},
		Client{},
		GqlAccessModel{},
		RedirectTarget{},
		RefreshToken{},
		RestAccessModel{},
		Session{},
		User{},
		Scope{},
	}
}
