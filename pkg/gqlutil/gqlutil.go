// Package gqlutil exposes GraphQL utility functionality
package gqlutil

import (
	"context"
	"strings"

	"github.com/machinebox/graphql"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

const introspectionQuery = `{
  __schema {
    types {
      name
      fields {
        name
        type {
          name
          kind
          ofType {
            name
            kind
            ofType {
              name
              kind
              ofType {
                name
                kind
              }
            }
          }
        }
      }
    }
  }
}`

// AccessModel GQL access model
type AccessModel struct {
	F []string               `json:"f"`
	M map[string]string      `json:"m"`
	P map[string]interface{} `json:"p,omitempty"`
}

type gqlIntroResponse struct {
	Schema gqlSchemaType `json:"__schema"`
}

type gqlSchemaType struct {
	Types []gqlType
}

type gqlType struct {
	Name   string
	Fields []gqlField
}

type gqlField struct {
	Name *string
	Kind string
	Type *ofType
}

type ofType struct {
	Name   *string
	Kind   string
	OfType *ofType
}

// Schema GQL Schema
type Schema struct {
	Types map[string]*Type
}

// Type GQL Type with fields
type Type struct {
	Fields map[string]*Field
}

// Field GQL Field
type Field struct {
	IsModel  bool
	TypeName string
}

// ErrGQLUtil is a global variable for error handling
var ErrGQLUtil = errors.New("something wrong in the gqlutil package")

// SchemaFetcher makes the FetchSchema function testable
type SchemaFetcher interface {
	FetchSchema(ctx context.Context, url string) (*Schema, error)
}

// GQLClient makes the Run function of the gql client testable
type GQLClient interface {
	// Run executes gql query on gqlClient
	Run(ctx context.Context, req *graphql.Request, resp interface{}) error
}

type schemaFetcher struct {
	createClient NewGraphQLClient
}

// NewGraphQLClient function to create a new GraphQL client
type NewGraphQLClient = func(string) GQLClient

// DefaultGraphQLClient default graphql client creation
func DefaultGraphQLClient(url string) GQLClient {
	return graphql.NewClient(url)
}

// NewSchemaFetcher creates a new SchemaFetcher for testing purposes
func NewSchemaFetcher(createClient NewGraphQLClient) SchemaFetcher {
	return &schemaFetcher{
		createClient: createClient,
	}
}

// FetchSchema fetches a GQL Schema at given location
func (fetcher schemaFetcher) FetchSchema(ctx context.Context, url string) (*Schema, error) {
	gqlResp := &gqlIntroResponse{}

	gqlClient := fetcher.createClient(url)

	err := gqlClient.Run(ctx, graphql.NewRequest(introspectionQuery), gqlResp)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch GQL scheme")
	}

	gqlSchema := gqlResp.Schema

	schema := Schema{Types: make(map[string]*Type)}
	for i := range gqlSchema.Types {
		tp := gqlSchema.Types[i]
		newType := Type{Fields: make(map[string]*Field)}
		for j := range tp.Fields {
			field := tp.Fields[j]
			fieldOfType := field.Type
			for fieldOfType.Name == nil {
				if fieldOfType.OfType == nil {
					return nil, errors.New("nil name without ofType in introspection query response")
				}
				fieldOfType = fieldOfType.OfType
			}
			fieldTypeName := *fieldOfType.Name
			if field.Name == nil {
				return nil, errors.New("field name is nil")
			}
			newType.Fields[*field.Name] = &Field{
				IsModel:  fieldOfType.Kind == "OBJECT" || fieldOfType.Kind == "INTERFACE",
				TypeName: fieldTypeName,
			}
		}
		schema.Types[tp.Name] = &newType
	}

	return &schema, nil
}

// ValidateQueryModel validates rights of an access model
func (s *Schema) ValidateQueryModel(queryModel map[string]*AccessModel) error { // nolint: funlen, gocognit, gocyclo
	isReferenced := make(map[string]bool)

	modelTypeNames := make(map[string]string)

	modelsToCheck := make(chan string, 1)

	rootModel, ok := queryModel["r"]
	if !ok {
		return errors.Wrap(ErrGQLUtil, "root model \"r\" not specified")
	}
	if len(rootModel.F) > 0 {
		return errors.Wrap(ErrGQLUtil, "root model cannot have fields")
	}

	modelTypeNames["r"] = "Query"
	modelsToCheck <- "r"

	modelsLeft := true
	for modelsLeft {
		select {
		// Keep processing models so long as they are left
		case modelRef := <-modelsToCheck:
			isReferenced[modelRef] = true
			model, ok := queryModel[modelRef]
			if !ok {
				return errors.Wrapf(ErrGQLUtil, "model %s is referenced but not specified", modelRef)
			}

			// Get type
			modelType, ok := s.Types[modelTypeNames[modelRef]]
			if !ok {
				return errors.Wrapf(ErrGQLUtil, "model %s is of type %s but said type is not specified", modelRef, modelTypeNames[modelRef])
			}

			// Check fields
			for _, fieldName := range model.F {
				gqlField, ok := modelType.Fields[fieldName]
				if !ok {
					return errors.Wrapf(ErrGQLUtil, "model %s of type %s does not have field %s", modelRef, modelTypeNames[modelRef], fieldName)
				}
				if gqlField.IsModel {
					return errors.Wrapf(ErrGQLUtil, "model %s has field %s, but %s is a model", modelRef, fieldName, fieldName)
				}
			}

			// Check submodels
			for submodelName, submodelReference := range model.M {
				modelField, ok := modelType.Fields[submodelName]
				if !ok {
					return errors.Wrapf(ErrGQLUtil, "model %s does not have submodel %s", modelRef, submodelName)
				}
				if !modelField.IsModel {
					return errors.Wrapf(ErrGQLUtil, "model %s has submodel %s, but %s is a field", modelRef, submodelName, submodelName)
				}

				reference, err := parseReference(submodelReference)
				if err != nil {
					return err
				}
				if existingModelTypeName, ok := modelTypeNames[reference]; ok {
					if existingModelTypeName != modelField.TypeName {
						return errors.Wrapf(ErrGQLUtil, "model %s is referenced in two places with different types: %s and %s",
							reference, existingModelTypeName, modelField.TypeName)
					}
				} else {
					modelTypeNames[reference] = modelField.TypeName
					modelsToCheck <- reference
				}
			}
		default:
			modelsLeft = false
		}
	}
	close(modelsToCheck)

	// Check whether all models have been referenced
	for k := range queryModel {
		if _, ok := isReferenced[k]; !ok {
			return errors.Wrapf(ErrGQLUtil, "model %s was specified but not referenced", k)
		}
	}

	return nil
}

func parseReference(modelReference string) (string, error) {
	if !strings.HasPrefix(modelReference, "#") {
		return "", errors.Wrapf(ErrGQLUtil, "model reference %s does not start with #", modelReference)
	}
	if !(len(modelReference) > 1) {
		return "", errors.Wrapf(ErrGQLUtil, "model reference %s is not valid", modelReference)
	}

	return modelReference[1:], nil
}
