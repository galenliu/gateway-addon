package addon

import json "github.com/json-iterator/go"

type Notification func(action *Action)

type Action struct {
	Name        string `json:"name"`
	AtType      string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ActionFunc  Notification

	DeviceId string `json:"deviceId"`
}

func NewAction() *Action {
	action := &Action{}
	return action
}

func (action *Action) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(action, "", " ")
}
