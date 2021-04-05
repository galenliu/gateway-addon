package addon

import (
	"addon/wot"
)

type Notification func(action *Action)

type Action struct {
	*wot.InteractionAffordance
	Name string `json:"name"`

	Input      *wot.DataSchema `json:"input,omitempty"`
	Output     *wot.DataSchema `json:"output,omitempty"`
	Safe       bool            `json:"safe"`
	Idempotent bool            `json:"idempotent"`
	DeviceId   string          `json:"deviceId"`
}

func NewActionFromString(data string) *Action {
	a := Action{}
	return &a
}

func NewAction() *Action {
	action := &Action{}
	return action
}

func (a *Action) AsDict() Map {
	return Map{
		"name": a.Name,
	}

}
