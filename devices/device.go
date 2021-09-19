package devices

import (
	"github.com/galenliu/gateway-addon/actions"
	"github.com/galenliu/gateway-addon/events"
	"github.com/galenliu/gateway-addon/properties"
	"github.com/galenliu/gateway-grpc"
	json "github.com/json-iterator/go"
)

type Device struct {
	ID          string `json:"id"`
	AtContext   string `json:"@context"`
	AtType      string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`

	Links               []*rpc.Link                     `json:"links,omitempty"`
	PinRequired         bool                            `json:"pinRequired"`
	CredentialsRequired bool                            `json:"credentialsRequired"`
	Pin                 Pin                             `json:"pin"`
	Properties          map[string]*properties.Property `json:"properties,omitempty"`
	Actions             map[string]*actions.Action      `json:"action,omitempty"`
	Events              map[string]*events.Event        `json:"events,omitempty"`
}

type Pin struct {
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
}

func NewDeviceFormMessage(dev *rpc.Device) *Device {
	device := &Device{
		ID:                  dev.Id,
		AtContext:           dev.AtContext,
		AtType:              dev.AtType,
		Title:               dev.Title,
		Description:         dev.Description,
		Links:               dev.Links,
		PinRequired:         dev.Pin.Required,
		CredentialsRequired: dev.CredentialsRequired,
		Pin: Pin{
			Required: dev.Pin.Required,
			Pattern:  dev.Pin.Pattern,
		},
	}
	if len(dev.Properties) > 0 {
		device.Properties = make(map[string]*properties.Property)
		for name, property := range dev.Properties {
			device.Properties[name] = properties.NewPropertyFormMessage(property)
		}
	}

	if len(dev.Events) > 0 {
		device.Events = make(map[string]*events.Event)
		for name, event := range dev.Events {
			device.Events[name] = events.NewEventFormMessage(event)
		}
	}

	if len(dev.Actions) > 0 {
		device.Actions = make(map[string]*actions.Action)
		for name, action := range dev.Actions {
			device.Actions[name] = actions.NewActionFormMessage(action)
		}
	}

	return device
}

func NewDeviceFormString(des string) *Device {
	var device Device
	err := json.UnmarshalFromString(des, &device)
	if err != nil {
		return nil
	}
	return &device
}

func (device *Device) GetID() string {
	return device.ID
}

func (device *Device) GetAtContext() string {
	return device.AtType
}

func (device *Device) GetAtType() string {
	return device.AtType
}

func (device *Device) GetTitle() string {
	return device.Title
}

func (device *Device) GetDescription() string {
	return device.Description
}

func (device *Device) AddProperty(property *properties.Property) {
	if device.Properties == nil {
		device.Properties = make(map[string]*properties.Property)
	}
	device.Properties[property.Name] = property
}
