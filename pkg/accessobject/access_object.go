// Package accessobject provides util functionality for objects defining access to graphql models
package accessobject

import (
	"encoding/json"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

// ErrIncorrectModelName incorrect model name
var ErrIncorrectModelName error = errors.New("incorrect model name - model references should look like \"#modelName\"")

// AccessModel access model defines how models and fields may be accessed in a graphql endpoint
type AccessModel struct {
	Fields     []string                           `json:"f,omitempty"`
	Models     map[string]*ReferenceOrAccessModel `json:"m,omitempty"`
	Parameters map[string]interface{}             `json:"p,omitempty"`
}

// ReferenceOrAccessModel struct containing either a reference or an access model
type ReferenceOrAccessModel struct {
	*AccessModel
	ref *string
}

type (
	// AccessObject is a map of keys and access models
	AccessObject map[string]*AccessModel
)

// UnmarshalJSON unmarshals data into a map of access models
func (a *AccessObject) UnmarshalJSON(data []byte) error {
	var ret map[string]*AccessModel

	if err := json.Unmarshal(data, &ret); err != nil {
		return errors.Wrap(err, "unable to unmarshal data into access object")
	}

	a.resolveModels(ret)

	*a = ret

	return nil
}

func (a *AccessObject) resolveModels(m map[string]*AccessModel) {
	models, references := findModelsAndReferences(m)

	for _, reference := range references {
		if reference.ref != nil && reference.AccessModel == nil {
			if f, ok := models[*reference.ref]; ok {
				reference.AccessModel = f
			}
		}
	}
}

func findModelsAndReferences(a map[string]*AccessModel) (map[string]*AccessModel, []*ReferenceOrAccessModel) {
	m := make(map[string]*AccessModel)
	var r []*ReferenceOrAccessModel

	for k, o := range a {
		m[k] = o

		for _, ref := range o.Models {
			if ref.AccessModel == nil && ref.ref != nil {
				r = append(r, ref)
			}
		}
	}

	return m, r
}

// GetModel retrieves model by key
func (a *AccessObject) GetModel(m string) *AccessModel {
	return (*a)[m]
}

// Root retrieves the root access model
func (a *AccessObject) Root() *AccessModel {
	return a.GetModel("r")
}

// UnmarshalJSON unmarshals byte array into ReferenceOrAccesModel struct
func (r *ReferenceOrAccessModel) UnmarshalJSON(data []byte) error {
	var check interface{}
	if err := json.Unmarshal(data, &check); err != nil {
		return errors.Wrap(err, "unable to unmarshal data into reference or access model")
	}

	var ret ReferenceOrAccessModel

	if ref, ok := check.(string); ok {
		if ref[0] != '#' {
			return errors.Wrap(ErrIncorrectModelName, "unable to unmarshal json into reference or access model")
		}

		modelReference := ref[1:]

		ret.ref = &modelReference
	} else {
		var am AccessModel
		if err := json.Unmarshal(data, &am); err != nil {
			return errors.Wrap(err, "unable to unmarshal data into access model")
		}
		ret.AccessModel = &am
	}

	*r = ret

	return nil
}
