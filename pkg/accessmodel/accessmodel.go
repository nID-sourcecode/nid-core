// Package accessmodel provides utility for parsing access models of various kinds
package accessmodel

import (
	"encoding/json"
	goErr "errors"

	"github.com/dgrijalva/jwt-go"

	"lab.weave.nl/nid/nid-core/pkg/accessobject"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

// Type is an enum containing the possible types of access models
type Type string

// Access model types
const (
	GQLType  Type = "GQL"
	RESTType Type = "REST"
)

// AccessModel is a generic access model wrapper
type AccessModel struct {
	Type            Type
	GQLAccessModel  GQLAccessModel
	RESTAccessModel RESTAccessModel
}

// GQLAccessModel is the type of access model for GQL scopes
type GQLAccessModel struct {
	Model accessobject.AccessObject `json:"m"`
	Path  string                    `json:"p"`
}

// RESTAccessModel is the type of access model for REST scopes
type RESTAccessModel struct {
	Path   string            `json:"p"`
	Query  map[string]string `json:"q"`
	Body   string            `json:"b"`
	Method string            `json:"m"`
}

// ErrInvalidAccessModelType is returned if the type value of an access model is unknown
var ErrInvalidAccessModelType = goErr.New("invalid access model type")

// UnmarshalJSON unmarshals an access model from json bytes
func (a *AccessModel) UnmarshalJSON(data []byte) error {
	var typeWrapper struct {
		Type Type `json:"t"`
	}

	err := json.Unmarshal(data, &typeWrapper)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal data into access model")
	}

	accessModel := AccessModel{
		Type: typeWrapper.Type,
	}
	switch typeWrapper.Type {
	case GQLType:
		err := json.Unmarshal(data, &accessModel.GQLAccessModel)
		if err != nil {
			return errors.Wrap(err, "unmarshalling gql access model")
		}
	case RESTType:
		err := json.Unmarshal(data, &accessModel.RESTAccessModel)
		if err != nil {
			return errors.Wrap(err, "umarshalling rest access model")
		}
	default:
		return ErrInvalidAccessModelType
	}

	*a = accessModel

	return nil
}

type claims struct {
	Scopes map[string]*AccessModel
}

func (*claims) Valid() error {
	return nil
}

// ExtractScopesFromJWT extracts scopes from JWT claims
func ExtractScopesFromJWT(token string) (map[string]*AccessModel, error) {
	var claims claims

	jwtParser := jwt.Parser{}
	if _, _, err := jwtParser.ParseUnverified(token, &claims); err != nil {
		return nil, errors.Wrap(err, "unable to parse unverified token")
	}

	return claims.Scopes, nil
}

// FilterByType filters scopes by their type
func FilterByType(scopes map[string]*AccessModel, scopeType Type) map[string]*AccessModel {
	return Filter(scopes, func(model *AccessModel) bool {
		return model.Type == scopeType
	})
}

// Filter filters scopes based on an arbitrary filtering condition
func Filter(scopes map[string]*AccessModel, condition func(*AccessModel) bool) map[string]*AccessModel {
	filteredScopes := make(map[string]*AccessModel)
	for k, v := range scopes {
		if condition(v) {
			filteredScopes[k] = v
		}
	}

	return filteredScopes
}
