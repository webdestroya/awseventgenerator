# awseventgenerator

Generates Go (golang) Structs and Validation code from AWS EventBridge JSON schemas.

# Requirements

* Go 1.20.3+

# Usage

Install

```console
$ go get -u github.com/webdestroya/awseventgenerator
```

# Example

This schema

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "title": "AwsExampleEventBridgeEvent",
  "definitions": {
    "TrafficLightChange": {
      "properties": {
        "new-state": {
          "$ref": "#/definitions/TrafficLightState"
        },
        "previous-state": {
          "$ref": "#/definitions/TrafficLightState"
        },
        "id": {
          "type": "string"
        }
      },
      "required": ["id", "new-state", "previous-state"],
      "type": "object"
    },
    "TrafficLightState": {
      "properties": {
        "blinking": {
          "type": "boolean"
        },
        "color": {
          "enum": ["RED", "YELLOW", "GREEN"],
          "type": "string"
        }
      },
      "required": ["color"],
      "type": "object"
    }
  },
  "properties": {
    "account": {
      "type": "string"
    },
    "detail": {
      "$ref": "#/definitions/TrafficLightChange"
    },
    "detail-type": {
      "type": "string"
    },
    "id": {
      "type": "string"
    },
    "region": {
      "type": "string"
    },
    "resources": {
      "items": {
        "type": "string"
      },
      "type": "array"
    },
    "source": {
      "type": "string"
    },
    "time": {
      "format": "date-time",
      "type": "string"
    },
    "version": {
      "type": "string"
    }
  },
  "required": [
    "detail-type",
    "resources",
    "id",
    "source",
    "time",
    "detail",
    "region",
    "version",
    "account"
  ],
  "type": "object",
  "x-amazon-events-detail-type": "Example EventBridge Event",
  "x-amazon-events-source": "aws.example"
}
```

generates

```go
package awsexample

import (
	"time"
)

const (
	AwsEventSource     = `aws.example`
	AwsEventDetailType = `Example EventBridge Event`
)

type ColorType string
const (
	ColorTypeRed    ColorType = "RED"
	ColorTypeYellow ColorType = "YELLOW"
	ColorTypeGreen  ColorType = "GREEN"
)
func (ColorType) Values() []ColorType {
	return []ColorType{
		"RED",
		"YELLOW",
		"GREEN",
	}
}

type AwsEvent struct {
	Account    string              `json:"account"`
	Detail     *TrafficLightChange `json:"detail"`
	DetailType string              `json:"detail-type"`
	Id         string              `json:"id"`
	Region     string              `json:"region"`
	Resources  []string            `json:"resources"`
	Source     string              `json:"source"`
	Time       time.Time           `json:"time"`
	Version    string              `json:"version"`
}

type TrafficLightChange struct {
	Id            string             `json:"id"`
	NewState      *TrafficLightState `json:"new-state"`
	PreviousState *TrafficLightState `json:"previous-state"`
}

type TrafficLightState struct {
	Blinking *bool     `json:"blinking,omitempty"`
	Color    ColorType `json:"color"`
}
```

To view more examples, run `make generate` and then look in the `internal/testcode/` directory.
