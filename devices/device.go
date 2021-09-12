package devices

import (
	"github.com/galenliu/gateway-addon/properties"
	json "github.com/json-iterator/go"
)

type Device struct {
	AdapterId string `json:"adapterId"`

	ID          string `json:"id"`
	AtContext   string `json:"@context"`
	AtType      string `json:"@type"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Links               []string `json:"links"`
	BaseHref            string   `json:"baseHref"`
	PinRequired         bool     `json:"pinRequired"`
	CredentialsRequired bool     `json:"credentialsRequired"`

	Properties map[string]*properties.Property `json:"properties"`
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

func (device *Device) GetName() string {
	return device.Name
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
