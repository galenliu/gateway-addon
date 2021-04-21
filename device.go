package addon

import (
	"fmt"
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

	//Properties map[string]*Property `json:"properties,omitempty"`
	Properties map[string]IProperty `json:"properties"`
	Actions    map[string]IAction   `json:"actions,omitempty"`
	Events     map[string]IEvent    `json:"events,omitempty"`

	Pin       PIN `json:"pin,omitempty"`
	username  string
	password  string
	AdapterId string `json:"adapterId"`

	adapter Owner
}

func NewDeivceFormString(data string) *Device {
	device := Device{}
	device.ID = gjson.Get(data, "id").String()
	device.Title = gjson.Get(data, "title").String()

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

	properties := gjson.Get(data, "properties").Map()
	if len(properties) > 1 {
		device.Properties = make(map[string]IProperty)
		for name, prop := range properties {
			p := NewPropertyFromString(prop.String())
			p.DeviceId = device.ID
			if p != nil {
				device.Properties[name] = p
			}
		}
	}

	actions := gjson.Get(data, "actions").Map()
	if len(actions) > 1 {
		device.Actions = make(map[string]IAction)
		for name, a := range properties {
			action := NewActionFromString(a.String())
			action.DeviceId = device.ID
			if action != nil {
				device.Actions[name] = action
			}
		}
	}

	events := gjson.Get(data, "events").Map()
	if len(events) > 1 {
		device.Actions = make(map[string]IAction)
		for name, e := range events {
			event := NewActionFromString(e.String())
			event.DeviceId = device.ID
			if event != nil {
				device.Events[name] = event
			}
		}
	}

	return &device
}

func NewDevice(id, title string) *Device {
	device := &Device{}
	device.ID = id
	device.Title = title
	device.AtType = make([]string, 0)
	device.AtContext = make([]string, 0)
	device.Properties = make(map[string]IProperty, 5)
	device.Actions = make(map[string]IAction, 1)
	device.Events = make(map[string]IEvent, 1)
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

func (device *Device) ToString() string {
	data, err := json.MarshalIndent(device, "", " ")
	if err != nil {
		return string(data)
	}
	return ""
}

//func (device *Device) MarshalWebThing() ([]byte, error) {
//
//}

func (device *Device) GetAdapterId() string {
	return device.AdapterId
}
