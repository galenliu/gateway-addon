package events

import rpc "github.com/galenliu/gateway-grpc"

type Event struct {
	Name string `json:"name"`
}

func NewEventFormMessage(event *rpc.Event) *Event {
	e := &Event{}
	return e
}
