package addon

//import json "github.com/json-iterator/go"

type Event struct {
	Name string
}

func (event *Event) GetEvent() *Event {
	return event
}
