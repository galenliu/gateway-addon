package addon

import "github.com/galenliu/gateway-addon/wot"

type Action struct {
	*wot.ActionAffordance
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
