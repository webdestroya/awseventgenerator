package testcode

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awseventgenerator"
	obj_vs_additional "github.com/webdestroya/awseventgenerator/internal/testcode/obj_vs_additional_gen"
)

func TestObjVsAdditionalGenerate(t *testing.T) {
	_, err := awseventgenerator.GenerateFromSchemaFile("../testdata/obj_vs_additional.json", &awseventgenerator.Config{
		GenerateEnums:           true,
		GenerateEnumValueMethod: true,
	})
	require.NoError(t, err)
}

func TestObjVsAdditional(t *testing.T) {

	objAll := map[string]interface{}{
		"strVal":   "somestring",
		"boolVal":  true,
		"nilVal":   nil,
		"floatVal": 123.123,
		"intVal":   1234,
		"arrVal":   []interface{}{123, 123.123, "strVal2", false, nil},
		"objVal": map[string]interface{}{
			"objk1": "something",
			"objk2": 123,
		},
	}

	t.Run("valid_everything", func(t *testing.T) {
		jsonMap := map[string]interface{}{
			"thing1": "val1",
			"thing2": map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
			},
			"thing3": objAll,
			"thing4": map[string]interface{}{
				"yar":      "blah",
				"somenum":  1234.123,
				"someprop": "propval",
				"foo":      "bar",
			},
			"thing5": objAll,
			"thing6": map[string]interface{}{
				"somenum":  1234.123,
				"someprop": "propval",
				"someobj": map[string]interface{}{
					"objk1": "something",
					"objk2": 123,
				},
				"nilVal":   nil,
				"floatVal": 123.123,
				"intVal":   1234,
			},
		}
		jsonData := jsonify(jsonMap)

		val := obj_vs_additional.Root{}
		err := json.Unmarshal([]byte(jsonData), &val)
		require.NoError(t, err, "Unmarshalling")

		marshData, err := json.Marshal(val)
		require.NoError(t, err, "marshalling")

		require.JSONEq(t, jsonData, string(marshData), "comparison")

		require.Equal(t, jsonMap["thing1"], *val.Thing1)

		// jmap2 := jsonMap["thing2"].(map[string]string)
		// require.Equal(t, jmap2["test"], *val.Thing2.)

		// jmap3 := jsonMap["thing3"].(map[string]interface{})

	})
}
