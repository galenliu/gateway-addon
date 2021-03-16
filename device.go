package addon

import (
	"fmt"
	json "github.com/json-iterator/go"
	//json "github.com/json-iterator/go"
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
	Required bool        `json:"required"`
	Pattern  interface{} `json:"pattern,omitempty"`
}

type Device struct {
	ID                  string   `json:"id"`
	AtContext           []string `json:"@context,omitempty"`
	Title               string   `json:"title,required"`
	AtType              []string `json:"@type"`
	Description         string   `json:"description,omitempty"`
	CredentialsRequired bool     `json:"credentialsRequired"`

	Properties map[string]*Property `json:"properties,omitempty"`
	Actions    map[string]*Action   `json:"actions,omitempty"`
	Events     map[string]*Event    `json:"events,omitempty"`

	Pin       PIN `json:"pin,omitempty"`
	username  string
	password  string
	AdapterId string `json:"adapterId"`
}

func NewDevice(id, title string) *Device {
	device := &Device{}
	device.ID = id
	device.Title = title
	device.AtType = make([]string, 0)
	device.AtContext = make([]string, 0)
	device.Properties = make(map[string]*Property, 5)
	device.Actions = make(map[string]*Action, 1)
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

func (device *Device) AddProperty(prop *Property) {
	if device.Properties == nil {
		device.Properties = make(map[string]*Property, 8)
	}
	prop.DeviceId = device.ID
	device.Properties[prop.Name] = prop
}

func (device *Device) AddAction(name string, a *Action) {
	if device.Actions == nil {
		device.Actions = make(map[string]*Action, 5)
	}
	device.Actions[name] = a
}

func (device *Device) AddEvent(name string, e *Event) {
	if device.Events == nil {
		device.Events = make(map[string]*Event, 8)
	}
	device.Events[name] = e
}

func (device *Device) AddTypes(types ...string) {
	for _, t := range types {
		device.AtType = append(device.AtType, t)
	}
}

func (device *Device) GetProperty(propertyName string) *Property {
	prop, ok := device.Properties[propertyName]
	if !ok {
		return nil
	}
	if prop.DeviceId == "" {
		prop.DeviceId = device.ID
	}
	return prop
}

func (device *Device) FindProperty(propertyName string) (*Property, error) {
	prop, ok := device.Properties[propertyName]
	if !ok {
		return nil, fmt.Errorf("can not found property(deivce:%s propertyName:%s)", device.ID, propertyName)
	}
	if prop.DeviceId == "" {
		prop.DeviceId = device.ID
	}
	return prop, nil
}

func (device *Device) SetProperty(propertyName string, value interface{}) (interface{}, error) {
	prop, ok := device.Properties[propertyName]
	if !ok {
		return nil, fmt.Errorf("properties(%s) not found", propertyName)
	}
	prop.UpdateValue(value)
	return prop.Value, nil
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

func (device *Device) ToJSON() string {
	data, err := json.MarshalIndent(device, "", " ")
	if err != nil {
		return string(data)
	}
	return ""
}
