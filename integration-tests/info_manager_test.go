//go:build integration || to || files
// +build integration to files

package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
	"lab.weave.nl/nid/nid-core/svc/info-manager/proto"
)

type InfoManagerTestSuite struct {
	BaseTestSuite
}

// func TestInfoManagerTestSuite(t *testing.T) {
// 	suite.Run(t, &InfoManagerTestSuite{})
// }

const script = `--define the new schema with a name
name = "personContactable"

--define the schema fields
schema = {
    answer = false,
}

input = {
    id = ""
}

--implement the resolve function
function resolve()
   local data = graphql("http://open-databron/gql", [[
        query($id: UUID!){
          user(id: $id) {
            contactDetails {
              phone
            }
          }
        }
    ]], { id = input.id })

    --if phoneNumber exists then person is contactable
    local phoneNumber = data.user.contactDetails[1].phone
    if phoneNumber == nil or phoneNumber == '' then
        schema["answer"] = false
    else
        schema["answer"] = true
    end

    return schema
end
`

func (s *InfoManagerTestSuite) TestScriptCRUD() {
	s.Run("Create", func() {
		createScriptMutation := gqlclient.NewRequest(`
			mutation {
				createScript(input: {
					name: "test",
					description: "test",
					status: DRAFT,
				}) {
					id
					description
					name
					status
				}
			}  
		`)

		createScriptMutationResponse := struct {
			CreateScript struct {
				ID          uuid.UUID `json:"id"`
				Name        string    `json:"name"`
				Description string    `json:"description"`
				Status      string    `json:"status"`
			}
		}{}

		err := s.infoManagerGQLClient.Post(context.Background(), createScriptMutation, &createScriptMutationResponse)
		s.Require().NoError(err)
		s.Require().NotNil(createScriptMutationResponse)
		s.Require().NotNil(createScriptMutationResponse.CreateScript.ID)
		s.Equal("test", createScriptMutationResponse.CreateScript.Name)
		s.Equal("test", createScriptMutationResponse.CreateScript.Description)
		s.Equal("DRAFT", createScriptMutationResponse.CreateScript.Status)

		random1, err := uuid.NewV4()
		s.Require().NoError(err)
		random2, err := uuid.NewV4()
		s.Require().NoError(err)

		s.Run("Get", func() {
		})

		s.Run("Update", func() {
			updateScriptMutation := gqlclient.NewRequest(fmt.Sprintf(`
				mutation {
					updateScript(
					id: "%s",	
					input: {
						name: "%s",
						description: "%s",
					}) {
						id
						description
						name
						status
					}
				}  
			`, createScriptMutationResponse.CreateScript.ID, random1, random2))

			updateScriptMutationResponse := struct {
				UpdateScript struct {
					ID          uuid.UUID `json:"id"`
					Name        string    `json:"name"`
					Description string    `json:"description"`
					Status      string    `json:"status"`
				}
			}{}

			err := s.infoManagerGQLClient.Post(context.Background(), updateScriptMutation, &updateScriptMutationResponse)
			s.Require().NoError(err)
			s.Require().NotNil(updateScriptMutationResponse)
			s.NotNil(updateScriptMutationResponse.UpdateScript.ID)
			s.Equal(random1.String(), updateScriptMutationResponse.UpdateScript.Name)
			s.Equal(random2.String(), updateScriptMutationResponse.UpdateScript.Description)
			s.Equal("DRAFT", updateScriptMutationResponse.UpdateScript.Status)
		})
	})
}

// Test a simple list scripts gql request
func (s *InfoManagerTestSuite) TestListScripts() {
	listScriptsQuery := gqlclient.NewRequest(`
		{
			scripts(limit: -1) {
			id
			name
			description
			status
			scriptSources {
				id
				signedUrl
				version
				changeDescription
				createdAt
				updatedAt
			}
			}
		}
	`)

	listScriptsResponse := struct {
		Scripts []struct {
			ID            uuid.UUID `json:"id"`
			Name          string    `json:"name"`
			Description   string    `json:"description"`
			Status        string    `json:"status"`
			ScriptSources []struct {
				ID                uuid.UUID `json:"id"`
				RawScript         string    `json:"rawScript"`
				SignedURL         string    `json:"signedUrl"`
				Version           string    `json:"version"`
				ChangeDescription string    `json:"changeDescription"`
				CreatedAt         time.Time `json:"createdAt"`
				UpdatedAt         time.Time `json:"updatedAt"`
			}
		}
	}{}

	err := s.infoManagerGQLClient.Post(context.Background(), listScriptsQuery, &listScriptsResponse)
	s.Require().NoError(err)
	s.Require().NotNil(listScriptsResponse)
	s.Require().True(len(listScriptsResponse.Scripts) > 1)
	s.NotNil(listScriptsResponse.Scripts[0].ID.String())
	s.NotNil(listScriptsResponse.Scripts[0].Name)
	s.NotNil(listScriptsResponse.Scripts[0].Description)
	s.NotNil(listScriptsResponse.Scripts[0].Status)

	s.Require().NotNil(listScriptsResponse.Scripts[0].ScriptSources)
	s.Require().Len(listScriptsResponse.Scripts[0].ScriptSources, 1)
	s.NotNil(listScriptsResponse.Scripts[0].ScriptSources[0].ID.String())
	s.NotNil(listScriptsResponse.Scripts[0].ScriptSources[0].ChangeDescription)
	s.NotNil(listScriptsResponse.Scripts[0].ScriptSources[0].Version)
}

// Test a more complex flow where we first create a script via gql and add a source via the grpc client
func (s *InfoManagerTestSuite) TestCreateWithLua() {
	s.Run("CreateScript", func() {
		createScriptMutation := gqlclient.NewRequest(`
			mutation {
				createScript(input: {
					name: "test",
					description: "test",
					status: DRAFT,
				}) {
					id
					description
					name
					status
				}
			}  
		`)

		createScriptMutationResponse := struct {
			CreateScript struct {
				ID          uuid.UUID `json:"id"`
				Name        string    `json:"name"`
				Description string    `json:"description"`
				Status      string    `json:"status"`
			}
		}{}

		err := s.infoManagerGQLClient.Post(context.Background(), createScriptMutation, &createScriptMutationResponse)
		s.Require().NoError(err)
		s.Require().NotNil(createScriptMutationResponse)
		s.Require().NotNil(createScriptMutationResponse.CreateScript.ID)
		s.Equal("test", createScriptMutationResponse.CreateScript.Name)
		s.Equal("test", createScriptMutationResponse.CreateScript.Description)
		s.Equal("DRAFT", createScriptMutationResponse.CreateScript.Status)

		s.Run("UploadLua", func() {
			scriptsUploadRequest := &proto.ScriptsUploadRequest{
				Script:            []byte(script),
				ScriptId:          createScriptMutationResponse.CreateScript.ID.String(),
				ChangeDescription: "Upload initial script",
			}

			_, err := s.infoManagerClient.ScriptsUpload(s.ctx, scriptsUploadRequest)
			s.Require().NoError(err)

			s.Run("GetSourcesForScript", func() {
				getScriptWithLua := gqlclient.NewRequest(fmt.Sprintf(`
					{
						script(id: "%s") {
						id
						name
						description
						status
						scriptSources {
						id
						rawScript
						signedUrl
						version
						changeDescription
						createdAt
						updatedAt
						}
					}
					}
				`, createScriptMutationResponse.CreateScript.ID))

				getScriptWithLuaResponse := struct {
					Script struct {
						ID            uuid.UUID `json:"id"`
						Name          string    `json:"name"`
						Description   string    `json:"description"`
						Status        string    `json:"status"`
						ScriptSources []struct {
							ID                uuid.UUID `json:"id"`
							RawScript         string    `json:"rawScript"`
							SignedURL         string    `json:"signedUrl"`
							Version           string    `json:"version"`
							ChangeDescription string    `json:"changeDescription"`
							CreatedAt         time.Time `json:"createdAt"`
							UpdatedAt         time.Time `json:"updatedAt"`
						}
					}
				}{}

				err := s.infoManagerGQLClient.Post(context.Background(), getScriptWithLua, &getScriptWithLuaResponse)
				s.Require().NoError(err)
				s.Require().NotNil(getScriptWithLuaResponse.Script)
				s.Equal(createScriptMutationResponse.CreateScript.ID, getScriptWithLuaResponse.Script.ID)
				s.Equal(createScriptMutationResponse.CreateScript.Name, getScriptWithLuaResponse.Script.Name)
				s.Equal(createScriptMutationResponse.CreateScript.Description, getScriptWithLuaResponse.Script.Description)
				s.Equal(createScriptMutationResponse.CreateScript.Status, getScriptWithLuaResponse.Script.Status)

				s.Require().NotNil(getScriptWithLuaResponse.Script.ScriptSources)
				s.Require().Len(getScriptWithLuaResponse.Script.ScriptSources, 1)
				s.Require().NotNil(getScriptWithLuaResponse.Script.ScriptSources[0].ID)
				s.Equal("Upload initial script", getScriptWithLuaResponse.Script.ScriptSources[0].ChangeDescription)

				s.Run("GetSourceWithLua", func() {
					getScriptSource := gqlclient.NewRequest(fmt.Sprintf(`
					{
						scriptSource(id: "%s") {
							id
							rawScript
							signedUrl
							version
							changeDescription
							createdAt
							updatedAt
						}
					}
					`, getScriptWithLuaResponse.Script.ScriptSources[0].ID))

					getScriptSourceResponse := struct {
						ScriptSource struct {
							ID                uuid.UUID `json:"id"`
							RawScript         string    `json:"rawScript"`
							SignedURL         string    `json:"signedUrl"`
							Version           string    `json:"version"`
							ChangeDescription string    `json:"changeDescription"`
							CreatedAt         time.Time `json:"createdAt"`
							UpdatedAt         time.Time `json:"updatedAt"`
						}
					}{}

					err := s.infoManagerGQLClient.Post(context.Background(), getScriptSource, &getScriptSourceResponse)
					s.Require().NoError(err)
					s.Require().NotNil(getScriptSourceResponse.ScriptSource)
					s.Equal(getScriptWithLuaResponse.Script.ScriptSources[0].ID, getScriptSourceResponse.ScriptSource.ID)
					s.Equal(getScriptWithLuaResponse.Script.ScriptSources[0].ChangeDescription, getScriptSourceResponse.ScriptSource.ChangeDescription)
					s.Equal(script, getScriptSourceResponse.ScriptSource.RawScript)
					s.NotEqual("", getScriptSourceResponse.ScriptSource.SignedURL)
				})
			})
		})
	})
}
