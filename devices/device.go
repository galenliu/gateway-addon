package devices

import (
	"github.com/galenliu/gateway-addon/properties"
	"github.com/galenliu/gateway-grpc"
	json "github.com/json-iterator/go"
)

type Device struct {
	AdapterId   string `json:"adapterId"`
	ID          string `json:"id"`
	AtContext   string `json:"@context"`
	AtType      string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Links               []*rpc.Link `json:"links"`
	BaseHref            string      `json:"baseHref"`
	PinRequired         bool        `json:"pinRequired"`
	CredentialsRequired bool        `json:"credentialsRequired"`
	Pin                 *rpc.DevicePin
	Properties          map[string]*properties.Property `json:"properties"`
}

func NewDeviceFormMessage(dev *rpc.Device) *Device {
	device := &Device{
		AdapterId:           "",
		ID:                  dev.Id,
		AtContext:           dev.AtContext,
		AtType:              dev.AtType,
		Title:               dev.Title,
		Description:         dev.Description,
		Links:               dev.Links,
		BaseHref:            dev.BaseHref,
		PinRequired:         dev.Pin.Required,
		CredentialsRequired: dev.CredentialsRequired,
		Pin:                 dev.Pin,
	}
	if len(dev.Properties) > 0 {
		device.Properties = make(map[string]*properties.Property)
	}
	for name, property := range dev.Properties {
		device.Properties[name] = properties.NewPropertyFormMessage(property)
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
