package addon

import "addon/wot"

//import json "github.com/json-iterator/go"

type Event struct {
	Name         string          `json:"name"`
	Subscription *wot.DataSchema `json:"subscription,omitempty"`
	Data         *wot.DataSchema `json:"data,omitempty"`
	Cancellation *wot.DataSchema `json:"cancellation,omitempty"`
	DeviceId     string          `json:"deviceId"`
}

func (event *Event) GetEvent() *Event {
	return event
}
