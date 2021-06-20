package addon

import (
	"github.com/galenliu/gateway/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
)

//import json "github.com/json-iterator/go"

type Event struct {
	Name         string                  `json:"name"`
	Subscription *data_schema.DataSchema `json:"subscription,omitempty"`
	Data         *data_schema.DataSchema `json:"data,omitempty"`
	Cancellation *data_schema.DataSchema `json:"cancellation,omitempty"`
	DeviceId     string                  `json:"deviceId"`
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
