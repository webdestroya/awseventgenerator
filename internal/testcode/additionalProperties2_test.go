package testcode

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	additionalProperties2 "github.com/webdestroya/awseventgenerator/internal/testcode/normal_gen/additionalproperties2_gen"
)

func TestMarshalUnmarshal(t *testing.T) {
	params := []struct {
		Name       string
		Strct      additionalProperties2.Root
		Validation func(t *testing.T, prop *additionalProperties2.Root)
	}{
		{
			Name: "Base Object",
			Strct: additionalProperties2.Root{
				Property1: ptr("test"),
			},
			Validation: func(t *testing.T, prop *additionalProperties2.Root) {
				if *prop.Property1 != "test" {
					t.Fatal("property1 != test")
				}
			},
		},
		{
			Name: "Property7",
			Strct: additionalProperties2.Root{
				Property7: &additionalProperties2.Property7{
					StreetNumber: ptr[float64](69),
					StreetName:   ptr("Elm St"),
					PoBox: &additionalProperties2.PoBox{
						Suburb: ptr("Smallville"),
					},
					AdditionalProperties: map[string]map[string]additionalProperties2.Anonymous1{
						"red": {
							"blue": {
								Color: ptr("green"),
								Conditions: []additionalProperties2.ConditionsItems{
									{Name: ptr("dry")},
								},
								Density: ptr(42.42),
							},
						},
						"orange": {},
					},
				},
			},
			Validation: func(t *testing.T, prop *additionalProperties2.Root) {

				if *prop.Property7.StreetNumber != 69 {
					t.Fatal("wrong value")
				}

				if len(prop.Property7.AdditionalProperties) != 2 {
					t.Fatal("not enough additionalProperties")
				}

				if *prop.Property7.AdditionalProperties["red"]["blue"].Color != "green" {
					t.Fatal("wrong nested value")
				}

				if *prop.Property7.AdditionalProperties["red"]["blue"].Density != 42.42 {
					t.Fatal("wrong nested value")
				}
			},
		},
	}

	for _, p := range params {
		if str, err := json.MarshalIndent(&p.Strct, "", "  "); err != nil {
			t.Fatal(err)
		} else {
			strct2 := &additionalProperties2.Root{}
			if err := json.Unmarshal(str, &strct2); err != nil {
				t.Fatal(err)
			}

			if reflect.DeepEqual(p.Strct, strct2) {
				log.Fatal("unmarshaled struct != given struct")
			}

			p.Validation(t, strct2)

			if str, err := json.MarshalIndent(&strct2, "", "  "); err != nil {
				t.Fatal(err)
			} else {
				strct3 := &additionalProperties2.Root{}
				if err := json.Unmarshal(str, &strct3); err != nil {
					t.Fatal(err)
				}

				if reflect.DeepEqual(p.Strct, strct3) {
					log.Fatal("unmarshaled struct != given struct")
				}

				p.Validation(t, strct3)

			}

		}
	}
}
