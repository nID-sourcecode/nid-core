package gqlutil

import (
	"context"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
)

type GQLClientMock struct {
	mock.Mock
}

// Run mocks the gql clients Run method
func (client *GQLClientMock) Run(ctx context.Context, req *graphql.Request, resp interface{}) error {
	args := client.Called(ctx, req, resp)

	return args.Error(0)
}

type SchemaFetcherTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	SchemaFetcher SchemaFetcher
	response      *gqlIntroResponse
	schema        *Schema
}

type QueryModelValidatorTestSuite struct {
	grpctesthelpers.GrpcTestSuite

	DefaultFetcher *SchemaFetcherTestSuite
	TestingFetcher *SchemaFetcherTestSuite
}

func (s *SchemaFetcherTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()

	gqlClientMock := &GQLClientMock{}
	gqlClientMock.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		response := args.Get(2).(*gqlIntroResponse)
		*response = *s.response
	})

	fetcher := NewSchemaFetcher(func(url string) GQLClient {
		return gqlClientMock
	})
	s.SchemaFetcher = fetcher
}

// Test the FetchSchema function with mocked values
// Creates a correct Schema which is also used in testing the QueryModelValidator
// Should return no error and a schema
func (s *SchemaFetcherTestSuite) TestFetchSchemaCorrect() {
	fieldName := "EpicField"
	ofTypeName := "EpicOfType"
	ofTypeTestName := "EpicTestType"
	ofTypeTestNested := ofType{
		Name:   &ofTypeName,
		Kind:   "OBJECT",
		OfType: nil,
	}
	ofTypeTest := ofType{
		Name:   nil,
		Kind:   "INTERFACE",
		OfType: &ofTypeTestNested,
	}
	ofTypeTestFields := ofType{
		Name:   &ofTypeTestName,
		Kind:   "EpicTestKind",
		OfType: nil,
	}
	var fields []gqlField
	fields = append(fields, gqlField{
		Name: &fieldName,
		Kind: "test",
		Type: &ofTypeTest,
	})
	fields = append(fields, gqlField{
		Name: &fieldName,
		Kind: "NoName",
		Type: &ofTypeTestNested,
	})
	id := "id"
	test := "test"
	var testFields []gqlField
	testFields = append(testFields, gqlField{
		Name: &id,
		Kind: "NotModel",
		Type: &ofTypeTestFields,
	})
	testFields = append(testFields, gqlField{
		Name: &test,
		Kind: "NotModel",
		Type: &ofTypeTestFields,
	})
	var types []gqlType
	types = append(types, gqlType{
		Name:   "Query",
		Fields: fields,
	}, gqlType{
		Name:   "EpicOfType",
		Fields: testFields,
	})

	s.response = &gqlIntroResponse{Schema: gqlSchemaType{types}}
	schema, err := s.SchemaFetcher.FetchSchema(s.Ctx, "url to fetch")

	field := Field{
		IsModel:  true,
		TypeName: "EpicOfType",
	}
	idField := Field{
		IsModel:  false,
		TypeName: "EpicTestType",
	}
	testField := Field{
		IsModel:  false,
		TypeName: "EpicTestType",
	}

	expected := Schema{Types: make(map[string]*Type)}
	expectedQueryType := Type{Fields: make(map[string]*Field)}
	expectedQueryType.Fields["EpicField"] = &field
	expectedTestType := Type{Fields: make(map[string]*Field)}
	expectedTestType.Fields["id"] = &idField
	expectedTestType.Fields["test"] = &testField
	// fix me
	expected.Types["Query"] = &expectedQueryType
	expected.Types["EpicOfType"] = &expectedTestType

	s.Require().NoError(err)
	s.Require().Equal(expected, *schema)
	s.schema = schema
}

// Tests what happens when Field name == nil
// Should return error
func (s *SchemaFetcherTestSuite) TestFetchSchemaFieldNameNil() {
	ofTypeTestName := "ofTypeTestName"
	ofTypeTest := ofType{
		Name:   &ofTypeTestName,
		Kind:   "OBJECT",
		OfType: nil,
	}
	var fields []gqlField
	fields = append(fields, gqlField{
		Name: nil,
		Kind: "test",
		Type: &ofTypeTest,
	})
	var types []gqlType
	types = append(types, gqlType{
		Name:   "schemaTypes",
		Fields: fields,
	})

	s.response = &gqlIntroResponse{Schema: gqlSchemaType{types}}
	_, err := s.SchemaFetcher.FetchSchema(s.Ctx, "url to fetch")

	s.Require().Error(err)
}

// TestFetchSchemaIncorrect tests what happens when both Oftype name and his nested Oftype are nil.
// Should return error
func (s *SchemaFetcherTestSuite) TestFetchSchemaIncorrect() {
	fieldName := "EpicField"
	ofTypeTest := ofType{
		Name:   nil,
		Kind:   "OBJECT",
		OfType: nil,
	}
	var fields []gqlField
	fields = append(fields, gqlField{
		Name: &fieldName,
		Kind: "test",
		Type: &ofTypeTest,
	})
	var types []gqlType
	types = append(types, gqlType{
		Name:   "schemaTypes",
		Fields: fields,
	})

	s.response = &gqlIntroResponse{Schema: gqlSchemaType{types}}
	_, err := s.SchemaFetcher.FetchSchema(s.Ctx, "url to fetch")

	s.Require().Error(err)
}

// TestFetchSchemaRunError tests the result of an error in gqlClient Run function
// Should return error
func (s *SchemaFetcherTestSuite) TestFetchSchemaRunError() {
	gqlClientMock := &GQLClientMock{}
	err := errors.Wrap(ErrGQLUtil, "Error in GQLClient.Run()")
	gqlClientMock.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(err)
	fetcher := NewSchemaFetcher(func(url string) GQLClient {
		return gqlClientMock
	})
	s.SchemaFetcher = fetcher

	_, err = s.SchemaFetcher.FetchSchema(s.Ctx, "url to fetch")
	s.Require().Error(err)
}

func (s *QueryModelValidatorTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()

	gqlClientMockDefault := &GQLClientMock{}
	gqlClientMockDefault.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		response := args.Get(2).(*gqlIntroResponse)
		*response = *s.DefaultFetcher.response
	})

	s.DefaultFetcher = &SchemaFetcherTestSuite{
		GrpcTestSuite: grpctesthelpers.GrpcTestSuite{},
		SchemaFetcher: NewSchemaFetcher(func(url string) GQLClient {
			return gqlClientMockDefault
		}),
		response: nil,
		schema:   nil,
	}
	// s.DefaultFetcher.SchemaFetcher = NewSchemaFetcher(gqlClientMock)
	s.DefaultFetcher.schema = s.DefaultFetcher.DefaultSchema()

	gqlClientMockTesting := &GQLClientMock{}
	gqlClientMockTesting.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		response := args.Get(2).(*gqlIntroResponse)
		*response = *s.TestingFetcher.response
	})
	// s.TestingFetcher.SchemaFetcher = NewSchemaFetcher(gqlClientMock)
	s.TestingFetcher = &SchemaFetcherTestSuite{
		GrpcTestSuite: grpctesthelpers.GrpcTestSuite{},
		SchemaFetcher: NewSchemaFetcher(func(url string) GQLClient {
			return gqlClientMockTesting
		}),
		response: nil,
		schema:   nil,
	}
	s.TestingFetcher.schema = s.TestingFetcher.SchemaForTesting()
}

// function to define the default schema for test cases
func (s *SchemaFetcherTestSuite) DefaultSchema() *Schema {
	rootTypeName := "rootType"
	rootFieldName := "rootField"
	rootType := ofType{
		Name:   &rootTypeName,
		Kind:   "OBJECT",
		OfType: nil,
	}
	testTypeName := "testType"
	testType := ofType{
		Name:   &testTypeName,
		Kind:   "NotModel",
		OfType: nil,
	}

	var fields []gqlField
	fields = append(fields, gqlField{
		Name: &rootFieldName,
		Kind: "",
		Type: &rootType,
	})

	id := "id"
	var subModelFields []gqlField
	subModelFields = append(subModelFields, gqlField{
		Name: &id,
		Kind: "NotModel",
		Type: &testType,
	})

	var types []gqlType
	types = append(types, gqlType{
		Name:   "Query",
		Fields: fields,
	}, gqlType{
		Name:   "rootType",
		Fields: subModelFields,
	})

	s.response = &gqlIntroResponse{Schema: gqlSchemaType{types}}
	schema, err := s.SchemaFetcher.FetchSchema(s.Ctx, "url to fetch")
	if err != nil {
		log.WithError(err).Fatal("schema could not be fetched")
	}

	return schema
}

// function to redefine the schema for FieldIsModelTest
func (s *SchemaFetcherTestSuite) SchemaForTesting() *Schema {
	rootTypeName := "rootType"
	rootFieldName := "rootField"
	fieldModelName := "fieldModel"
	rootType := ofType{
		Name:   &rootTypeName,
		Kind:   "OBJECT",
		OfType: nil,
	}
	testTypeName := "testType"
	testType := ofType{
		Name:   &testTypeName,
		Kind:   "NotModel",
		OfType: nil,
	}
	fieldModel := ofType{
		Name:   &fieldModelName,
		Kind:   "NotModel",
		OfType: nil,
	}

	var fields []gqlField
	fields = append(fields, gqlField{
		Name: &rootFieldName,
		Kind: "",
		Type: &rootType,
	})

	id := "id"
	fieldIsModel := "fieldIsModel"
	var subModelFields []gqlField
	subModelFields = append(subModelFields, gqlField{
		Name: &id,
		Kind: "NotModel",
		Type: &testType,
	}, gqlField{
		Name: &fieldIsModel,
		Kind: "OBJECT",
		Type: &rootType,
	}, gqlField{
		Name: &fieldModelName,
		Kind: "NotModel",
		Type: &fieldModel,
	})

	var types []gqlType
	types = append(types, gqlType{
		Name:   "Query",
		Fields: fields,
	}, gqlType{
		Name:   "rootType",
		Fields: subModelFields,
	})

	s.response = &gqlIntroResponse{Schema: gqlSchemaType{types}}
	schema, err := s.SchemaFetcher.FetchSchema(s.Ctx, "url to fetch")
	if err != nil {
		log.WithError(err).Fatal("schema could not be fetched")
	}

	return schema
}

const accessModelKey = "#testValue"

// TestValidateQueryModelCorrect tests a correct model being validated.
// Should return no error
func (s *QueryModelValidatorTestSuite) TestValidateQueryModelCorrect() {
	accessModels := make(map[string]*AccessModel)
	rootModel := make(map[string]string)
	var value []string
	value = append(value, "id")
	rootModel["rootField"] = accessModelKey
	root := AccessModel{
		F: nil,
		M: rootModel,
		P: nil,
	}
	testValue := AccessModel{
		F: value,
		M: nil,
		P: nil,
	}
	accessModels["r"] = &root
	accessModels["testValue"] = &testValue

	err := s.DefaultFetcher.schema.ValidateQueryModel(accessModels)
	s.Require().NoError(err)
}

// TestValidateQueryRootModelNotSpecified checks when the rootmodel is not present
// should return error
func (s *QueryModelValidatorTestSuite) TestValidateQueryModelRootModelNotSpecified() {
	accessModels := make(map[string]*AccessModel)
	rootModel := make(map[string]string)
	var value []string
	value = append(value, "id")
	rootModel["rootField"] = accessModelKey
	testValue := AccessModel{
		F: value,
		M: nil,
		P: nil,
	}
	accessModels["testValue"] = &testValue

	err := s.DefaultFetcher.schema.ValidateQueryModel(accessModels)
	s.Require().Error(err)
}

// TestValidateQueryModelRootModelWithFields tests when the rootmodel has field, which it should not have
// Should return error
func (s *QueryModelValidatorTestSuite) TestValidateQueryModelRootModelWithFields() {
	accessModels := make(map[string]*AccessModel)
	rootModel := make(map[string]string)
	var rootModelField []string
	rootModelField = append(rootModelField, "field")
	var value []string
	value = append(value, "id")
	rootModel["rootField"] = accessModelKey
	root := AccessModel{
		F: rootModelField,
		M: rootModel,
		P: nil,
	}
	testValue := AccessModel{
		F: value,
		M: nil,
		P: nil,
	}
	accessModels["r"] = &root
	accessModels["testValue"] = &testValue

	err := s.DefaultFetcher.schema.ValidateQueryModel(accessModels)
	s.Require().Error(err)
}

// TestValidateQueryModelNotSpecified tests when a model is referenced but not specified
// Should return error
func (s *QueryModelValidatorTestSuite) TestValidateQueryModelNotSpecified() {
	accessModels := make(map[string]*AccessModel)
	rootModel := make(map[string]string)
	rootModel["rootField"] = accessModelKey
	root := AccessModel{
		F: nil,
		M: rootModel,
		P: nil,
	}
	accessModels["r"] = &root

	err := s.DefaultFetcher.schema.ValidateQueryModel(accessModels)
	s.Require().Error(err)
}

// TestValidateQueryModelWrongField
// Should return error
func (s *QueryModelValidatorTestSuite) TestValidateQueryModelWrongField() {
	accessModels := make(map[string]*AccessModel)
	rootModel := make(map[string]string)
	var value []string
	value = append(value, "id", "wrongField")
	rootModel["rootField"] = accessModelKey
	root := AccessModel{
		F: nil,
		M: rootModel,
		P: nil,
	}
	testValue := AccessModel{
		F: value,
		M: nil,
		P: nil,
	}
	accessModels["r"] = &root
	accessModels["testValue"] = &testValue

	err := s.DefaultFetcher.schema.ValidateQueryModel(accessModels)
	s.Require().Error(err)
}

// TestValidateQueryModelWrongSubmodel tests when the wrong submodel is specified
// Should return error
func (s *QueryModelValidatorTestSuite) TestValidateQueryModelWrongSubmodel() {
	accessModels := make(map[string]*AccessModel)
	rootModel := make(map[string]string)
	var value []string
	value = append(value, "id")
	rootModel["failedField"] = accessModelKey
	root := AccessModel{
		F: nil,
		M: rootModel,
		P: nil,
	}
	testValue := AccessModel{
		F: value,
		M: nil,
		P: nil,
	}
	accessModels["r"] = &root
	accessModels["testValue"] = &testValue

	err := s.DefaultFetcher.schema.ValidateQueryModel(accessModels)
	s.Require().Error(err)
}

// TestValidateQueryModelFieldIsModel tests when a given field is actually a model
// Should return error
func (s *QueryModelValidatorTestSuite) TestValidateQueryModelFieldIsModel() {
	accessModels := make(map[string]*AccessModel)
	rootModel := make(map[string]string)
	var value []string
	value = append(value, "id", "fieldIsModel")
	rootModel["rootField"] = accessModelKey
	root := AccessModel{
		F: nil,
		M: rootModel,
		P: nil,
	}
	testValue := AccessModel{
		F: value,
		M: nil,
		P: nil,
	}
	accessModels["r"] = &root
	accessModels["testValue"] = &testValue
	err := s.TestingFetcher.schema.ValidateQueryModel(accessModels)
	s.Require().Error(err)
}

// TestValidateQueryModelSubmodelIsField test when a submodel is actually a field
// should return error
func (s *QueryModelValidatorTestSuite) TestValidateQueryModelSubmodelIsField() {
	accessModels := make(map[string]*AccessModel)
	rootModel := make(map[string]string)
	fieldModel := make(map[string]string)
	var value []string
	value = append(value, "id")
	rootModel["rootField"] = accessModelKey
	fieldModel["fieldModel"] = "#fieldModel"
	root := AccessModel{
		F: nil,
		M: rootModel,
		P: nil,
	}
	testValue := AccessModel{
		F: value,
		M: fieldModel,
		P: nil,
	}
	fieldModelValue := AccessModel{
		F: nil,
		M: nil,
		P: nil,
	}
	accessModels["r"] = &root
	accessModels["testValue"] = &testValue
	accessModels["fieldModel"] = &fieldModelValue

	err := s.TestingFetcher.schema.ValidateQueryModel(accessModels)
	s.Require().Error(err)
}

func TestSchemaFetcherTestSuite(t *testing.T) {
	suite.Run(t, &SchemaFetcherTestSuite{})
	suite.Run(t, &QueryModelValidatorTestSuite{})
}
