package addon

import (
	"fmt"
	"github.com/galenliu/gateway-addon/wot"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"log"
)

type ChangeFunc func(property *Property, newValue, oldValue interface{})
type GetFunc func() interface{}

type Owner interface {
	Send(int, map[string]interface{})
}

type Property struct {
	*wot.PropertyAffordance

	Name string `json:"name"`

	DeviceId string `json:"deviceId,omitempty"`

	device            Owner
	updateOnSameValue bool

	valueChangeFuncs []ChangeFunc
	valueGetFunc     GetFunc
	verbose          bool
}

func NewPropertyFromString(description string) *Property {
	var prop Property

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
	return &prop
}

func NewProperty(typ string) *Property {
	prop := &Property{
		valueChangeFuncs: make([]ChangeFunc, 0),
	}
	prop.SetAtType(typ)
	prop.verbose = true
	return prop
}

func (p *Property) OnValueUpdate(fn ChangeFunc) {
	p.valueChangeFuncs = append(p.valueChangeFuncs, fn)
}

func (p *Property) OnValueGet(fn GetFunc) {
	p.valueGetFunc = fn
}

func (p *Property) GetValue() interface{} {
	return p.getValue()
}

func (p *Property) getValue() interface{} {
	if p.valueGetFunc != nil {
		p.UpdateValue(p.valueGetFunc())
	}
	return p.Value
}

func (p *Property) SetCachedValueAndNotify(value interface{}) {
	p.UpdateValue(value)
	data := make(map[string]interface{})
	data["property"] = p.ToString()
	p.device.Send(DevicePropertyChangedNotification, data)
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

//func (prop *Property) MarshalJSON() ([]byte, error) {
//	return json.MarshalIndent(prop, "", " ")
//}

func (p *Property) ToString() string {
	str, err := json.MarshalToString(p)
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}
	return str
}

func (p *Property) SetOwner(owner Owner) {
	p.device = owner
}

func (p *Property) AsDict() []byte {
	m := Map{
		"name":        p.Name,
		"value":       p.Value,
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

func (p *Property) MarshalJson() []byte {
	data, err := json.Marshal(p)
	if err == nil {
		return data
	}
	return nil
}
