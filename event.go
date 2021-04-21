package addon

import (
	"github.com/galenliu/gateway-addon/wot"
	json "github.com/json-iterator/go"
)

//import json "github.com/json-iterator/go"

type Event struct {
	Name         string          `json:"name"`
	Subscription *wot.DataSchema `json:"subscription,omitempty"`
	Data         *wot.DataSchema `json:"data,omitempty"`
	Cancellation *wot.DataSchema `json:"cancellation,omitempty"`
	DeviceId     string          `json:"deviceId"`
}

func (e *Event) GetEvent() *Event {
	return e
}

func (e *Event) MarshalJson() []byte {
	data, err := json.Marshal(e)
	if err == nil {
		return data
	}
	return nil
}
