// Code generated by awseventgenerator/internal/generators/testcode. DO NOT EDIT.

package testsuitegenerated

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	issue6 "github.com/webdestroya/awseventgenerator/internal/testcode/issue6_gen"
)

func TestGenerated_issue6(t *testing.T) {

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

	t.Run("structs", func(t *testing.T) {
		t.Run("Attributes", func(t *testing.T) {
			genStruct := &issue6.Attributes{
				Hostname: &strVal,
				Location: &strVal,
				OperatingSystem: &issue6.OperatingSystem{
					Family:   &strVal,
					Name:     &strVal,
					Revision: &strVal,
					Version:  &strVal,
				},
				SerialNumber: &strVal,
			}
			t.Run("json", func(t *testing.T) {
				jsonOut, err := json.Marshal(genStruct)
				require.NoError(t, err)

				unmarObj := &issue6.Attributes{}
				require.NoError(t, json.Unmarshal(jsonOut, unmarObj))

				jsonOut2, err := json.Marshal(unmarObj)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonOut), string(jsonOut2))

				var jsearch interface{}
				require.NoError(t, json.Unmarshal(jsonOut, &jsearch))
				requireJmesMatch(t, jsearch, `"hostname"`, strVal, "Hostname")
				requireJmesMatch(t, jsearch, `"location"`, strVal, "Location")
				require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length("operating_system")`).(float64), 1.0)
				requireJmesMatch(t, jsearch, `"serial_number"`, strVal, "SerialNumber")

			})
			t.Run("fields", func(t *testing.T) {
				require.Equal(t, strVal, *issue6.Attributes{Hostname: &strVal}.Hostname)
				require.Equal(t, strVal, *issue6.Attributes{Location: &strVal}.Location)
				require.NotNil(t, genStruct.OperatingSystem) // Lazily Tested: issue6.Attributes.OperatingSystem == *OperatingSystem
				require.Equal(t, strVal, *issue6.Attributes{SerialNumber: &strVal}.SerialNumber)
			})
		})

		t.Run("LinksItems", func(t *testing.T) {
			genStruct := &issue6.LinksItems{
				AssetId:      &strVal,
				Description:  &strVal,
				Relationship: &strVal,
			}
			t.Run("json", func(t *testing.T) {
				jsonOut, err := json.Marshal(genStruct)
				require.NoError(t, err)

				unmarObj := &issue6.LinksItems{}
				require.NoError(t, json.Unmarshal(jsonOut, unmarObj))

				jsonOut2, err := json.Marshal(unmarObj)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonOut), string(jsonOut2))

				var jsearch interface{}
				require.NoError(t, json.Unmarshal(jsonOut, &jsearch))
				requireJmesMatch(t, jsearch, `"asset_id"`, strVal, "AssetId")
				requireJmesMatch(t, jsearch, `"description"`, strVal, "Description")
				requireJmesMatch(t, jsearch, `"relationship"`, strVal, "Relationship")

			})
			t.Run("fields", func(t *testing.T) {
				require.Equal(t, strVal, *issue6.LinksItems{AssetId: &strVal}.AssetId)
				require.Equal(t, strVal, *issue6.LinksItems{Description: &strVal}.Description)
				require.Equal(t, strVal, *issue6.LinksItems{Relationship: &strVal}.Relationship)
			})
		})

		t.Run("OperatingSystem", func(t *testing.T) {
			genStruct := &issue6.OperatingSystem{
				Family:   &strVal,
				Name:     &strVal,
				Revision: &strVal,
				Version:  &strVal,
			}
			t.Run("json", func(t *testing.T) {
				jsonOut, err := json.Marshal(genStruct)
				require.NoError(t, err)

				unmarObj := &issue6.OperatingSystem{}
				require.NoError(t, json.Unmarshal(jsonOut, unmarObj))

				jsonOut2, err := json.Marshal(unmarObj)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonOut), string(jsonOut2))

				var jsearch interface{}
				require.NoError(t, json.Unmarshal(jsonOut, &jsearch))
				requireJmesMatch(t, jsearch, `"family"`, strVal, "Family")
				requireJmesMatch(t, jsearch, `"name"`, strVal, "Name")
				requireJmesMatch(t, jsearch, `"revision"`, strVal, "Revision")
				requireJmesMatch(t, jsearch, `"version"`, strVal, "Version")

			})
			t.Run("fields", func(t *testing.T) {
				require.Equal(t, strVal, *issue6.OperatingSystem{Family: &strVal}.Family)
				require.Equal(t, strVal, *issue6.OperatingSystem{Name: &strVal}.Name)
				require.Equal(t, strVal, *issue6.OperatingSystem{Revision: &strVal}.Revision)
				require.Equal(t, strVal, *issue6.OperatingSystem{Version: &strVal}.Version)
			})
		})

		t.Run("Root", func(t *testing.T) {
			genStruct := &issue6.Root{
				Attributes: &issue6.Attributes{
					Hostname: &strVal,
					Location: &strVal,
					OperatingSystem: &issue6.OperatingSystem{
						Family:   &strVal,
						Name:     &strVal,
						Revision: &strVal,
						Version:  &strVal,
					},
					SerialNumber: &strVal,
				},
				Links: []issue6.LinksItems{{
					AssetId:      &strVal,
					Description:  &strVal,
					Relationship: &strVal,
				}},
				Tags: []string{strVal},
			}
			t.Run("json", func(t *testing.T) {
				jsonOut, err := json.Marshal(genStruct)
				require.NoError(t, err)

				unmarObj := &issue6.Root{}
				require.NoError(t, json.Unmarshal(jsonOut, unmarObj))

				jsonOut2, err := json.Marshal(unmarObj)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonOut), string(jsonOut2))

				var jsearch interface{}
				require.NoError(t, json.Unmarshal(jsonOut, &jsearch))
				require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length("attributes")`).(float64), 1.0)
				require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length("links")`).(float64), 1.0)
				require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length("tags")`).(float64), 1.0)

			})
			t.Run("fields", func(t *testing.T) {
				require.NotNil(t, genStruct.Attributes) // Lazily Tested: issue6.Root.Attributes == *Attributes
				require.NotNil(t, genStruct.Links)      // Lazily Tested: issue6.Root.Links == []LinksItems
				require.Contains(t, issue6.Root{Tags: []string{strVal}}.Tags, strVal)
			})
		})

	})

}
