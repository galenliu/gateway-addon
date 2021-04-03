package addon

//import json "github.com/json-iterator/go"

type Event struct {
	Name string `json:"name"`

	DeviceId string `json:"deviceId"`
}

func (event *Event) GetEvent() *Event {
	return event
}
