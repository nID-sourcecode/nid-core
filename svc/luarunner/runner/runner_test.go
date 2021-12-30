package runner

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	lua "github.com/yuin/gopher-lua"

	"lab.weave.nl/nid/nid-core/pkg/luautil"
)

func TestToLuaValueInteger(t *testing.T) {
	state := lua.NewState()

	result := luautil.ToLuaValue(state, float64(5))

	if result.String() != "5" {
		t.Error(fmt.Sprintf("result should be 5, was %s", result.String()))
	}

	if result.Type() != lua.LTNumber {
		t.Error(fmt.Sprintf("type should be LTNumber, was %s", result.Type()))
	}
}

func TestToLuaValueMap(t *testing.T) {
	state := lua.NewState()

	goValue := map[string]interface{}{
		"apples": "bananas",
		"pie":    5,
	}

	result := luautil.ToLuaValue(state, goValue)

	if result.Type() != lua.LTTable {
		t.Error(fmt.Sprintf("type should be LTTable, was %s", result.Type()))
	}
}

func TestScriptLoggingFailUnder18(t *testing.T) {
	runner := &LuaRunner{}

	file, err := ioutil.ReadFile("./gql-small-test.lua")
	if err != nil {
		t.Error(err)
	}

	err = runner.RunScript(context.Background(), string(file), map[string]interface{}{
		"geboorteDatum": "2018-12-23T00:00:00Z",
	})

	if err == nil {
		t.Error("Run script did not fail")
	}
}

func TestScriptLoggingOver18(t *testing.T) {
	runner := &LuaRunner{}

	file, err := ioutil.ReadFile("./gql-small-test.lua")
	if err != nil {
		t.Error(err)
	}

	err = runner.RunScript(context.Background(), string(file), map[string]interface{}{
		"geboorteDatum": "1997-07-23T00:00:00Z",
	})

	if err != nil {
		t.Error(err)
	}
}

func BenchmarkToLuaValueMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		state := lua.NewState()
		goValue := map[string]interface{}{
			"data": map[string]interface{}{
				"client": map[string]interface{}{
					"geheimeClient": false,
					"director": map[string]string{
						"name":     "Zorgkantoor: DaVinci",
						"endpoint": "https://gqlzkdavinci.staging.n-id.network/gql",
					},
					"person": map[string]string{
						"geboorteDatum":           "1991-02-01T00:00:00Z",
						"geslacht":                "Vrouwelijk",
						"geslachtsnaam":           "Dam",
						"voorletters":             "T.M.",
						"voorvoegselGeslachtnaam": "van",
					},
				},
			},
		}
		b.StartTimer()
		luautil.ToLuaValue(state, goValue)
	}
}
