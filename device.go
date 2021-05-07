package addon

import (
	"fmt"
	"github.com/galenliu/gateway-addon/wot"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

const (
	Alarm                    = "Alarm"
	AirQualitySensor         = "AirQualitySensor"
	BarometricPressureSensor = "BarometricPressureSensor"
	BinarySensor             = "BinarySensor"
	Camera                   = "Camera"
	ColorControl             = "ColorControl"
	ColorSensor              = "ColorSensor"
	DoorSensor               = "DoorSensor"
	EnergyMonitor            = "EnergyMonitor"
	HumiditySensor           = "HumiditySensor"
	LeakSensor               = "LeakSensor"
	Light                    = "Light"
	Lock                     = "Lock"
	MotionSensor             = "MotionSensor"
	MultiLevelSensor         = "MultiLevelSensor"
	MultiLevelSwitch         = "MultiLevelSwitch"
	OnOffSwitch              = "OnOffSwitch"
	SmartPlug                = "SmartPlug"
	SmokeSensor              = "SmokeSensor"
	TemperatureSensor        = "TemperatureSensor"
	Thermostat               = "Thermostat"
	VideoCamera              = "VideoCamera"
	Context                  = "https://webthings.io/schemas"
)

type PIN struct {
	Required bool        `json:"required,omitempty"`
	Pattern  interface{} `json:"pattern,omitempty"`
}

type AdapterProxy interface {
	Send(int, map[string]interface{})
}

type Device struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	AtContext           []string `json:"@context,omitempty"`
	Title               string   `json:"title"`
	Titles              []string `json:"titles,omitempty"`
	AtType              []string `json:"@type"`
	Description         string   `json:"description,omitempty"`
	CredentialsRequired bool     `json:"credentialsRequired,omitempty"`

	//Properties map[string]*Property `json:"properties,omitempty"`
	Properties map[string]IProperty `json:"properties"`
	Actions    map[string]IAction   `json:"actions,omitempty"`
	Events     map[string]IEvent    `json:"events,omitempty"`

	Pin       PIN `json:"pin,omitempty"`
	username  string
	password  string
	AdapterId string `json:"adapterId"`

	Forms []wot.Form `json:"forms,omitempty"`

	adapter AdapterProxy
}

func NewDeviceFormString(data string, adapter AdapterProxy) *Device {

	id := gjson.Get(data, "id").String()
	if id == "" {
		return nil
	}
	title := gjson.Get(data, "title").String()
	if title == "" {
		title = id
	}
	device := NewDevice(id, title, adapter)

	if gjson.Get(data, "@context").IsArray() {
		for _, c := range gjson.Get(data, "@context").Array() {
			device.AtContext = append(device.AtContext, c.String())
		}
	} else {
		if gjson.Get(data, "@context").String() != "" {
			device.AtContext = append(device.AtContext, gjson.Get(data, "@context").String())
		}
	}

	if gjson.Get(data, "@type").IsArray() {
		for _, c := range gjson.Get(data, "@type").Array() {
			device.AtType = append(device.AtType, c.String())
		}
	} else {
		if gjson.Get(data, "@type").String() != "" {
			device.AtType = append(device.AtType, gjson.Get(data, "@type").String())
		}
	}
	device.CredentialsRequired = gjson.Get(data, "credentialsRequired").Bool()

	var pin PIN
	if gjson.Get(data, "pin").Exists() {
		pin.Pattern = gjson.Get(data, "pin.pattern").Value()
		pin.Required = gjson.Get(data, "pin.required").Bool()
		device.Pin = pin
	}

	if gjson.Get(data, "properties").Exists() {
		properties := gjson.Get(data, "properties").Map()
		if len(properties) > 1 {

			for name, prop := range properties {
				p := NewPropertyFromString(prop.String(), device)
				p.DeviceId = device.ID
				if p != nil {
					device.Properties[name] = p
				}
			}
		}

	}

	if gjson.Get(data, "actions").Exists() {
		actions := gjson.Get(data, "actions").Map()
		if len(actions) > 1 {

			for name, a := range actions {
				action := NewActionFromString(a.String())
				action.DeviceId = device.ID
				if action != nil {
					device.Actions[name] = action
				}
			}
		}

	}

	if gjson.Get(data, "events").Exists() {
		events := gjson.Get(data, "events").Map()
		if len(events) > 1 {

			for name, e := range events {
				event := NewActionFromString(e.String())
				event.DeviceId = device.ID
				if event != nil {
					device.Events[name] = event
				}
			}
		}

	}

	return device
}

func NewDevice(id, title string, adp AdapterProxy) *Device {
	device := &Device{}
	device.ID = id
	device.Title = title
	device.AtType = make([]string, 0)
	device.AtContext = make([]string, 0)
	device.Properties = make(map[string]IProperty)
	device.Actions = make(map[string]IAction)
	device.Events = make(map[string]IEvent)
	if adp != nil {
		device.adapter = adp
	}
	return device
}

func (device *Device) GetTitle() string {
	return device.Title
}

func (device *Device) SetDescription(dsc string) {
	device.Description = dsc
}

func (device *Device) GetDescription() string {
	return device.Description
}

func (device *Device) SetTitle(title string) {
	device.Title = title
}

func (device *Device) AddProperty(prop IProperty) {
	if device.Properties == nil {
		device.Properties = make(map[string]IProperty, 8)
	}
	device.Properties[prop.GetName()] = prop
}

func (device *Device) AddAction(name string, a IAction) {
	device.Actions[name] = a
}

func (device *Device) AddEvent(name string, e IEvent) {
	device.Events[name] = e
}

func (device *Device) AddTypes(types ...string) {
	for _, t := range types {
		device.AtType = append(device.AtType, t)
	}
}

func (device *Device) GetProperty(propertyName string) IProperty {
	prop, ok := device.Properties[propertyName]
	if !ok {
		return nil
	}
	return prop
}

func (device *Device) Send(mt int, data map[string]interface{}) {
	data[Did] = device.GetID()
	device.adapter.Send(mt, data)
}

func (device *Device) GetID() string {
	return device.ID
}

func (device *Device) SetCredentials(username, password string) error {
	device.username = username
	device.password = password
	return nil
}

func (device *Device) SetPin(pin interface{}) error {
	if device.Pin.Required == false {
		return fmt.Errorf("devices pin not required")
	}
	device.Pin.Pattern = pin
	return nil
}

func (device *Device) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(device, "", " ")
}

func (device *Device) AsDict() Map {
	m := Map{
		"id":                  device.ID,
		"title":               device.Title,
		"@context":            device.AtContext,
		"@type":               device.AtType,
		"description":         device.Description,
		"forms":               device.Forms,
		"pin":                 device.Pin,
		"credentialsRequired": device.CredentialsRequired,
		"properties":          device.mapPropertiesToDict(),
		"events":              device.mapEventsDictFromFunction(),
		"actions":             device.mapActionsDictFromFunction(),
	}
	return m
}

func (device *Device) mapPropertiesToDict() Map {
	m := make(Map)
	for name, p := range device.Properties {
		m[name] = p
	}
	return m
}

func (device *Device) mapActionsDictFromFunction() Map {
	m := make(Map)
	for name, a := range device.Actions {
		m[name] = a
	}
	return m
}

func (device *Device) mapEventsDictFromFunction() Map {
	m := make(Map)
	for name, a := range device.Events {
		m[name] = a
	}
	return m
}

func (device *Device) GetAdapterId() string {
	return device.AdapterId
}
