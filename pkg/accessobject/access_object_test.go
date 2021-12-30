package accessobject

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AccessObjectTestSuite struct {
	suite.Suite
}

func TestAccessObjectTestSuite(t *testing.T) {
	suite.Run(t, new(AccessObjectTestSuite))
}

func (suite *AccessObjectTestSuite) TestUnmarshal() {
	var ao map[string]*AccessObject

	err := json.Unmarshal([]byte(`{
        "wlz": {
            "r": {
                "m": {
                    "client": "#C"
                }
            },
            "C": {
                "p": {
                    "orderBy": "bsn",
      				"complexParam": { "with": "children", "and": "stuff" },
      				"filter": null
                },
                "f": [
                    "geheimeClient"
                ],
                "m": {
                    "director": "#D",
                    "person": "#P"
                }
            },
            "D": {
                "f": ["name", "endpoint"]
            },
            "P": {
                "f": ["geboorteDatum", "geslacht", "geslachtsnaam", "voorletters", "voorvoegselGeslachtnaam"]
            }
        },
        "wallet": {
            "r": {
                "m": {
                    "wallet": "#W"
                }
            },
            "W": {
                "p": {
                    "orderBy": "bsn",
      				"complexParam": { "with": "children", "and": "stuff" },
      				"filter": null
                },
                "f": [
                    "phone"
                ]
            }
        }
    }`), &ao)
	suite.Require().NoError(err)

	foundKeys := make([]string, 0, len(ao))
	for k := range ao {
		foundKeys = append(foundKeys, k)
	}

	expectedKeys := []string{"wlz", "wallet"}
	sort.Strings(foundKeys)
	sort.Strings(expectedKeys)
	assert.Equal(suite.T(), expectedKeys, foundKeys, "Multiple access objects in map")

	assert.Nil(suite.T(), ao["wlz"].Root().Fields, "Fields should not be set for a root without fields")
	assert.Nil(suite.T(), ao["wlz"].Root().Parameters, "Parameters should not be set for a root without parameters")

	assert.NotNil(suite.T(), ao["wlz"].Root().Models, "Models should be set for a root with models")

	// Check if reference string is set
	assert.Equal(suite.T(), "C", *ao["wlz"].Root().Models["client"].ref, "The reference strings should be parsed (without the # prefix)")

	// Check if reference is resolved
	assert.Equal(suite.T(), ao["wlz"].GetModel("C"), ao["wlz"].Root().Models["client"].AccessModel, "The reference should be resolved")
}
