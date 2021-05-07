package addon

import (
	"github.com/galenliu/gateway-addon/wot"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"log"
)

type ChangeFunc func(property *Property, newValue, oldValue interface{})
type GetFunc func() interface{}

type DeviceProxy interface {
	Send(int, map[string]interface{})
}

type Property struct {
	*wot.PropertyAffordance

	Name     string `json:"name"`
	DeviceId string `json:"deviceId,omitempty"`

	device            DeviceProxy
	updateOnSameValue bool

	valueChangeFuncs []ChangeFunc
	valueGetFunc     GetFunc
	verbose          bool
}

func NewPropertyFromString(description string, proxy DeviceProxy) *Property {

	name := gjson.Get(description, "name").String()
	deviceId := gjson.Get(description, "deviceId").String()

	prop := NewProperty(name, deviceId, proxy, false)

	prop.PropertyAffordance = wot.NewPropertyAffordanceFromString(description)
	if gjson.Get(description, "name").Exists() {
		prop.Name = gjson.Get(description, "name").String()
	}

	if gjson.Get(description, "deviceId").Exists() {
		prop.DeviceId = gjson.Get(description, "deviceId").String()
	}

	if gjson.Get(description, "value").Exists() {
		if gjson.Get(description, "value").Value() != nil {
			prop.SetCachedValue(gjson.Get(description, "value").Value())
		} else {
			prop.SetCachedValue(prop.GetDefaultValue())
		}
	} else {
		prop.SetCachedValue(prop.GetDefaultValue())
	}

	prop.valueChangeFuncs = make([]ChangeFunc, 0)

	var onChanged = func(prop *Property, new interface{}, old interface{}) {
		if proxy != nil {
			data := make(map[string]interface{})
			data["property"] = prop.AsDict()
			proxy.Send(DevicePropertyChangedNotification, data)
		}
	}
	prop.valueChangeFuncs = append(prop.valueChangeFuncs, onChanged)
	return prop
}

func NewProperty(name, deviceId string, device DeviceProxy, verbose bool) *Property {
	prop := &Property{
		valueChangeFuncs: make([]ChangeFunc, 0),
	}
	prop.Name = name
	prop.DeviceId = deviceId
	if device != nil {
		prop.device = device
	}
	prop.verbose = verbose
	return prop
}

func (p *Property) OnValueUpdate(fn ChangeFunc) {
	p.valueChangeFuncs = append(p.valueChangeFuncs, fn)
}

func (p *Property) OnValueGet(fn GetFunc) {
	p.valueGetFunc = fn
}

// GetValue 获取Device Property Value

func (p *Property) GetValue() interface{} {
	return p.getValue()
}

//回调valueGetFunc，向设备发送GetValue请求。
func (p *Property) getValue() interface{} {
	if p.valueGetFunc != nil {
		p.UpdateValue(p.valueGetFunc())
	}
	return p.Value
}

func (p *Property) SetCachedValueAndNotify(value interface{}) {
	p.UpdateValue(value)
}

func (p *Property) UpdateValue(value interface{}) {

	if p.Value == value && !p.updateOnSameValue {
		return
	}
	if p.IsReadOnly() {
		return
	}
	oldValue := p.Value
	p.SetCachedValue(value)
	p.onValueUpdate(p.valueChangeFuncs, value, oldValue)
}

func (p *Property) onValueUpdate(funcs []ChangeFunc, newValue, oldValue interface{}) {
	for _, fn := range funcs {
		fn(p, newValue, oldValue)
	}
}

func (p *Property) DoPropertyChanged(d string) {
	title := gjson.Get(d, "title").String()
	if title != "" && p.Title != title {
		p.Title = title
	}
	p.SetType(gjson.Get(d, "type").String())
	p.SetAtType(gjson.Get(d, "@type").String())
	value := gjson.Get(d, "value").Value()
	if value != nil && p.Value != value {
		p.SetCachedValue(value)
	}
}

func (p *Property) UpdateProperty(d string) {
	title := gjson.Get(d, "title").String()
	if title != "" && p.Title != title {
		p.Title = title
	}
	p.SetType(gjson.Get(d, "type").String())
	p.SetAtType(gjson.Get(d, "@type").String())
	value := gjson.Get(d, "value").Value()
	if value != nil && p.Value != value {
		p.SetCachedValue(value)
	}
}

func (p *Property) SetValue(newValue interface{}) {
	if p.verbose {
		log.Printf("property(%s) set value not imp", p.GetName())
	}
}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) SetName(name string) {
	p.Name = name
}

func (p *Property) GetAtType() string {
	return p.AtType
}

func (p *Property) SetDeviceProxy(device IDevice) {
	p.device = device
}

func (p *Property) AsDict() []byte {
	m := Map{
		"name":        p.Name,
		"value":       p.GetValue(),
		"title":       p.Title,
		"type":        p.GetType(),
		"@type":       p.AtType,
		"description": p.Description,
		"readOnly":    p.IsReadOnly(),
		"forms":       p.Forms,
		"deviceId":    p.DeviceId,
	}
	data, err := json.MarshalIndent(&m, "", "  ")
	if err != nil {
		return nil
	}
	return data
}

func (p *Property) MarshalJSON() ([]byte, error) {
	m := Map{
		"name":        p.Name,
		"value":       p.GetValue(),
		"title":       p.Title,
		"type":        p.GetType(),
		"@type":       p.AtType,
		"description": p.Description,
		"readOnly":    p.IsReadOnly(),
		"forms":       p.Forms,
		"deviceId":    p.DeviceId,
	}
	data, err := json.MarshalIndent(&m, "", "  ")
	if err != nil {
		return nil, err
	}
	return data, nil
}
