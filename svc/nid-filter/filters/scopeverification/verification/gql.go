// Package verification provides various kinds of Verfiers, which verify whether a request fits a certain kind of scope.
package verification

import (
	"encoding/json"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"

	"github.com/nID-sourcecode/nid-core/pkg/accessmodel"
	"github.com/nID-sourcecode/nid-core/pkg/accessobject"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

type gqlRequestBody struct {
	Query     string
	Variables map[string]interface{}
}

// GQLVerifier verifies that the request matches any combination of the GQL scopes in the token.
type GQLVerifier struct{}

// Verify verifies that the request matches any combination of the GQL scopes in the token.
func (*GQLVerifier) Verify(req *Request) error {
	gqlScopes := accessmodel.FilterByType(req.Scopes, accessmodel.GQLType)

	scopes := accessmodel.Filter(gqlScopes, func(model *accessmodel.AccessModel) bool {
		return model.GQLAccessModel.Path == req.Path
	})

	if len(scopes) == 0 {
		return errors.Errorf("%w: no gql scopes found", ErrNotValid)
	}

	gqlRequestBody := gqlRequestBody{}

	switch req.Method {
	case "POST":
		if err := json.Unmarshal([]byte(req.Body), &gqlRequestBody); err != nil {
			return errors.Errorf("%w: parsing body: %v", ErrBadRequest, err)
		}
	case "GET":
		gqlRequestBody.Query = req.Query.Get("query")
		gqlRequestBody.Variables = make(map[string]interface{})
		variablesJSON := req.Query.Get("variables")
		err := json.Unmarshal([]byte(variablesJSON), &(gqlRequestBody.Variables))
		if err != nil {
			return errors.Errorf("%w: parsing variables: %v", ErrBadRequest, err)
		}
	default:
		return errors.Errorf("%w: method %s is not supported, use POST or GET", ErrBadRequest, req.Method)
	}

	query, err := parser.ParseQuery(&ast.Source{Input: gqlRequestBody.Query})
	if err != nil {
		return errors.Errorf("%w: parsing GQL query: %v", ErrBadRequest, err)
	}

	var rootScopes []*accessobject.ReferenceOrAccessModel
	for _, scope := range scopes {
		rootScopes = append(rootScopes, &accessobject.ReferenceOrAccessModel{AccessModel: scope.GQLAccessModel.Model.Root()})
	}

	for _, operation := range query.Operations {
		for _, sel := range operation.SelectionSet {
			selection := sel.(*ast.Field)

			// Skip introspection queries
			if strings.HasPrefix(selection.Name, "__") {
				continue
			}

			err := doesQueryMatchScope(gqlRequestBody.Variables, selection, "", rootScopes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func doesQueryMatchScope(vars map[string]interface{}, queryField *ast.Field, prefix string, scopeModels []*accessobject.ReferenceOrAccessModel) error {
	// Is queryField a field? I don't know.. is it?
	if len(queryField.SelectionSet) == 0 {
		for _, scopeModel := range scopeModels {
			// If so, check whether it is allowed
			if queryField.Name == "__typename" || stringInSlice(queryField.Name, scopeModel.Fields) {
				return nil
			}
		}

		return ErrNotValid
	}
	// queryField is a model
	// Find nested scope models within the current ones that match the queryField
	var childScopeModels []*accessobject.ReferenceOrAccessModel
	childScopeModelsFound := false
	for _, scopeModel := range scopeModels {
		if childModel, ok := scopeModel.Models[queryField.Name]; ok {
			childScopeModelsFound = true
			if verifyFilters(vars, queryField, childModel) == nil {
				childScopeModels = append(childScopeModels, childModel)
			}
		}
	}
	if len(childScopeModels) == 0 && childScopeModelsFound {
		return ErrNotValid
	}

	if len(childScopeModels) > 0 { // Continue checking
		for _, queryChildField := range queryField.SelectionSet {
			if err := doesQueryMatchScope(vars, queryChildField.(*ast.Field), prefix+queryField.Name+".", childScopeModels); err != nil {
				return err
			}
		}

		return nil
	}

	return ErrNotValid
}

var errInvalidFilter = errors.New("invalid filter")

func verifyFilters(vars map[string]interface{}, field *ast.Field, scopeModel *accessobject.ReferenceOrAccessModel) error {
	// Map to keep track of which required parameters have been specified in the query
	paramFound := make(map[string]bool)

	// Check whether specified parameters match their respective required parameter
	for _, queryParam := range field.Arguments {
		paramValue, err := queryParam.Value.Value(vars)
		if err != nil {
			return errors.Wrapf(err, "getting parameter value for %s", queryParam.Name)
		}
		paramShouldBe, ok := scopeModel.Parameters[queryParam.Name]
		if !ok {
			return errors.Errorf("%w: no access to parameter %s", errInvalidFilter, queryParam.Name)
		}
		equal, err := JSONEquals(paramValue, paramShouldBe)
		if err != nil {
			return errors.Wrap(err, "comparing parameter values")
		}
		if !equal {
			return errors.Errorf("%w: parameter value for %s was %+v but should be %+v", errInvalidFilter,
				queryParam.Name, paramValue, paramShouldBe)
		}
		paramFound[queryParam.Name] = true
	}

	// Check whether all required parameters have been specified
	for k := range scopeModel.Parameters {
		if _, ok := paramFound[k]; !ok {
			if scopeModel.Parameters[k] != nil { // Except for nil ones, which need not be set
				return errors.Errorf("%w: parameter value for %s is required but not set", errInvalidFilter, k)
			}
		}
	}

	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
