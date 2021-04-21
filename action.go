package addon

import (
	"github.com/galenliu/gateway-addon/wot"
	json "github.com/json-iterator/go"
)

type Action struct {
	*wot.ActionAffordance

	Name     string `json:"name"`
	DeviceId string `json:"deviceId"`
}

func NewActionFromString(data string) *Action {
	var a Action
	err := json.UnmarshalFromString(data, &a)
	if err != nil {
		return nil
	}
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

func (a *Action) MarshalJson() []byte {
	data, err := json.Marshal(a)
	if err == nil {
		return data
	}
	return nil
}
