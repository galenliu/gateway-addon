package wot

type Event struct {
	Name         string      `json:"name"`
	Subscription *DataSchema `json:"subscription,omitempty"`
	Data         *DataSchema `json:"data,omitempty"`
	Cancellation *DataSchema `json:"cancellation,omitempty"`
}

func (event *Event) GetEvent() *Event {
	return event
}
