// Code generated by awseventgenerator/internal/generators/testcode. DO NOT EDIT.

package testsuitegenerated

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	obj_vs_additional "github.com/webdestroya/awseventgenerator/internal/testcode/obj_vs_additional_gen"
)

func TestGenerated_obj_vs_additional(t *testing.T) {

	strVal := "someString"
	floatVal := float64(1232.1424)
	intVal := int64(1232)
	timeVal := time.Now().UTC()
	trueVal := true
	anyVal := struct {
		Thing string `json:"thinger"`
	}{Thing: "anywayanyday"}

	require.IsType(t, *new(string), strVal)
	require.IsType(t, *new(float64), floatVal)
	require.IsType(t, *new(int64), intVal)
	require.IsType(t, *new(time.Time), timeVal)
	require.IsType(t, *new(bool), trueVal)
	_ = anyVal

	t.Run("aliases", func(t *testing.T) {
		require.IsType(t, *new(map[string]obj_vs_additional.AddEmptyWithObjsItem), *new(obj_vs_additional.AddEmptyWithObjs))
		require.IsType(t, *new(interface{}), *new(obj_vs_additional.AddEmptyWithObjsItem))
		require.IsType(t, *new(map[string]string), *new(obj_vs_additional.AddEmptyWithStrings))
		require.IsType(t, *new(interface{}), *new(obj_vs_additional.ExtraWithObjsItem))
		require.IsType(t, *new(interface{}), *new(obj_vs_additional.PlainObj))
		require.IsType(t, *new(interface{}), *new(obj_vs_additional.Someobj))
	})

	t.Run("structs", func(t *testing.T) {
		t.Run("ExtraWithObjs", func(t *testing.T) {
			genStruct := &obj_vs_additional.ExtraWithObjs{
				AdditionalProperties: map[string]obj_vs_additional.ExtraWithObjsItem{
					strVal: anyVal,
				},
				Somenum:  &floatVal,
				Someobj:  anyVal,
				Someprop: &strVal,
			}
			t.Run("json", func(t *testing.T) {
				jsonOut, err := json.Marshal(genStruct)
				require.NoError(t, err)

				unmarObj := &obj_vs_additional.ExtraWithObjs{}
				require.NoError(t, json.Unmarshal(jsonOut, unmarObj))

				jsonOut2, err := json.Marshal(unmarObj)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonOut), string(jsonOut2))

				var jsearch interface{}
				require.NoError(t, json.Unmarshal(jsonOut, &jsearch))
				requireJmesMatch(t, jsearch, `"somenum"`, floatVal, "Somenum")
				require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length("someobj")`).(float64), 1.0)
				requireJmesMatch(t, jsearch, `"someprop"`, strVal, "Someprop")

			})
			t.Run("fields", func(t *testing.T) {
				require.NotNil(t, genStruct.AdditionalProperties) // Lazily Tested: obj_vs_additional.ExtraWithObjs.AdditionalProperties == map[string]ExtraWithObjsItem
				require.Equal(t, floatVal, *obj_vs_additional.ExtraWithObjs{Somenum: &floatVal}.Somenum)
				require.NotNil(t, genStruct.Someobj) // Lazily Tested: obj_vs_additional.ExtraWithObjs.Someobj == Someobj
				require.Equal(t, strVal, *obj_vs_additional.ExtraWithObjs{Someprop: &strVal}.Someprop)
			})
		})

		t.Run("ExtraWithStrings", func(t *testing.T) {
			genStruct := &obj_vs_additional.ExtraWithStrings{
				AdditionalProperties: map[string]string{
					strVal: strVal,
				},
				Somenum:  &floatVal,
				Someprop: &strVal,
			}
			t.Run("json", func(t *testing.T) {
				jsonOut, err := json.Marshal(genStruct)
				require.NoError(t, err)

				unmarObj := &obj_vs_additional.ExtraWithStrings{}
				require.NoError(t, json.Unmarshal(jsonOut, unmarObj))

				jsonOut2, err := json.Marshal(unmarObj)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonOut), string(jsonOut2))

				var jsearch interface{}
				require.NoError(t, json.Unmarshal(jsonOut, &jsearch))
				requireJmesMatch(t, jsearch, `"somenum"`, floatVal, "Somenum")
				requireJmesMatch(t, jsearch, `"someprop"`, strVal, "Someprop")

			})
			t.Run("fields", func(t *testing.T) {
				require.NotNil(t, genStruct.AdditionalProperties) // Lazily Tested: obj_vs_additional.ExtraWithStrings.AdditionalProperties == map[string]string
				require.Equal(t, floatVal, *obj_vs_additional.ExtraWithStrings{Somenum: &floatVal}.Somenum)
				require.Equal(t, strVal, *obj_vs_additional.ExtraWithStrings{Someprop: &strVal}.Someprop)
			})
		})

		t.Run("Root", func(t *testing.T) {
			genStruct := &obj_vs_additional.Root{
				Thing1: &strVal,
				Thing2: map[string]string{
					strVal: strVal,
				},
				Thing3: anyVal,
				Thing4: &obj_vs_additional.ExtraWithStrings{
					AdditionalProperties: map[string]string{
						strVal: strVal,
					},
					Somenum:  &floatVal,
					Someprop: &strVal,
				},
				Thing5: map[string]obj_vs_additional.AddEmptyWithObjsItem{
					strVal: anyVal,
				},
				Thing6: &obj_vs_additional.ExtraWithObjs{
					AdditionalProperties: map[string]obj_vs_additional.ExtraWithObjsItem{
						strVal: anyVal,
					},
					Somenum:  &floatVal,
					Someobj:  anyVal,
					Someprop: &strVal,
				},
			}
			t.Run("json", func(t *testing.T) {
				jsonOut, err := json.Marshal(genStruct)
				require.NoError(t, err)

				unmarObj := &obj_vs_additional.Root{}
				require.NoError(t, json.Unmarshal(jsonOut, unmarObj))

				jsonOut2, err := json.Marshal(unmarObj)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonOut), string(jsonOut2))

				var jsearch interface{}
				require.NoError(t, json.Unmarshal(jsonOut, &jsearch))
				requireJmesMatch(t, jsearch, `"thing1"`, strVal, "Thing1")
				require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length("thing2")`).(float64), 1.0)
				require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length("thing3")`).(float64), 1.0)
				require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length("thing4")`).(float64), 1.0)
				require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length("thing5")`).(float64), 1.0)
				require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length("thing6")`).(float64), 1.0)

			})
			t.Run("fields", func(t *testing.T) {
				require.Equal(t, strVal, *obj_vs_additional.Root{Thing1: &strVal}.Thing1)
				require.NotNil(t, genStruct.Thing2) // Lazily Tested: obj_vs_additional.Root.Thing2 == AddEmptyWithStrings
				require.NotNil(t, genStruct.Thing3) // Lazily Tested: obj_vs_additional.Root.Thing3 == PlainObj
				require.NotNil(t, genStruct.Thing4) // Lazily Tested: obj_vs_additional.Root.Thing4 == *ExtraWithStrings
				require.NotNil(t, genStruct.Thing5) // Lazily Tested: obj_vs_additional.Root.Thing5 == AddEmptyWithObjs
				require.NotNil(t, genStruct.Thing6) // Lazily Tested: obj_vs_additional.Root.Thing6 == *ExtraWithObjs
			})
		})

	})

}
