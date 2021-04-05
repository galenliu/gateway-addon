package addon

import (
	"addon/wot"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/xiam/to"
	"log"
)

type ChangeFunc func(property *Property, newValue, oldValue interface{})
type GetFunc func() interface{}

type Owner interface {
	Send(int, map[string]interface{})
}

type Property struct {
	*wot.DataSchema
	Name       string        `json:"name"`
	Value      interface{}   `json:"value"`
	Unit       string        `json:"unit,omitempty"`
	Minimum    interface{}   `json:"minimum,omitempty"`
	Maximum    interface{}   `json:"maximum,omitempty"`
	MultipleOf int           `json:"multipleOf,omitempty"`
	Enum       []interface{} `json:"enum,omitempty"`

	DeviceId string `json:"deviceId,omitempty"`

	replyChan chan Map

	device            Owner
	updateOnSameValue bool

	valueChangeFuncs []ChangeFunc
	valueGetFunc     GetFunc
	verbose          bool
}

func NewPropertyFromString(description string) *Property {
	var prop Property
	prop.replyChan = make(chan Map)
	json.UnmarshalFromString(description, &prop)
	if prop.Type == TypeNumber {
		if gjson.Get(description, "minimum").Exists() || gjson.Get(description, "minimum").Exists() {
			schema := wot.NumberSchema{}
			if gjson.Get(description, "minimum").Exists() {
				schema.Minimum = gjson.Get(description, "minimum").Float()
			}
			if gjson.Get(description, "maximum").Exists() {
				schema.Minimum = gjson.Get(description, "maximum").Float()
			}
			if gjson.Get(description, "exclusiveMinimum").Exists() {
				schema.ExclusiveMinimum = gjson.Get(description, "exclusiveMinimum").Float()
			}
			if gjson.Get(description, "exclusiveMaximum").Exists() {
				schema.ExclusiveMaximum = gjson.Get(description, "exclusiveMaximum").Float()
			}
			if gjson.Get(description, "multipleOf").Exists() {
				schema.MultipleOf = gjson.Get(description, "multipleOf").Float()
			}
			prop.Schema = schema

		}
	}
	if prop.Type == TypeInteger {
		if gjson.Get(description, "minimum").Exists() || gjson.Get(description, "minimum").Exists() {
			schema := wot.IntegerSchema{}
			if gjson.Get(description, "minimum").Exists() {
				schema.Minimum = gjson.Get(description, "minimum").Int()
			}
			if gjson.Get(description, "maximum").Exists() {
				schema.Minimum = gjson.Get(description, "maximum").Int()
			}
			if gjson.Get(description, "exclusiveMinimum").Exists() {
				schema.ExclusiveMinimum = gjson.Get(description, "exclusiveMinimum").Int()
			}
			if gjson.Get(description, "exclusiveMaximum").Exists() {
				schema.ExclusiveMaximum = gjson.Get(description, "exclusiveMaximum").Int()
			}
			if gjson.Get(description, "multipleOf").Exists() {
				schema.MultipleOf = gjson.Get(description, "multipleOf").Int()
			}
			prop.Schema = schema

		}
	}
	if prop.Type == TypeString {
		if gjson.Get(description, "minLength").Exists() || gjson.Get(description, "maxLength").Exists() {
			schema := wot.StringSchema{}
			if gjson.Get(description, "minLength").Exists() {
				schema.MinLength = gjson.Get(description, "minLength").Int()
			}
			if gjson.Get(description, "maximum").Exists() {
				schema.MaxLength = gjson.Get(description, "maxLength").Int()
			}
			prop.Schema = schema
		}
	}
	return &prop
}

func NewProperty(typ string) *Property {
	prop := &Property{
		valueChangeFuncs: make([]ChangeFunc, 0),
	}
	prop.DataSchema.AtType = typ
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
	value = p.convert(value)
	switch p.Type {
	case TypeNumber:
		value = p.clampFloat(value.(float64))
	case TypeInteger:
		value = p.clampInt(value.(int))
	}
	if p.Value == value && !p.updateOnSameValue {
		return
	}
	if p.ReadOnly {
		return
	}
	oldValue := p.Value
	p.Value = value
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
	value := gjson.Get(d, "value").Value()
	if value != nil && p.Value != value {
		p.Value = value
	}
	select {
	case p.replyChan <- p.AsDict():
		return
	default:
		return
	}
}

func (p *Property) convert(v interface{}) interface{} {
	switch p.Type {
	case TypeNumber:
		return to.Float64(v)
	case TypeInteger:
		return int(to.Uint64(v))
	case TypeBoolean:
		return to.Bool(v)
	default:
		return v
	}
}

func (p *Property) clampFloat(value float64) interface{} {
	min, minOK := p.Minimum.(float64)
	max, maxOK := p.Maximum.(float64)
	if maxOK == true && value > max {
		value = max
	} else if minOK == true && value < min {
		value = min
	}
	return value
}

func (p *Property) clampInt(value int) interface{} {
	min, minOK := p.Minimum.(int)
	max, maxOK := p.Maximum.(int)
	if maxOK == true && value > max {
		value = max
	} else if minOK == true && value < min {
		value = min
	}
	return value
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

func (p *Property) GetType() string {
	return p.Type
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

func (p *Property) AsDict() Map {
	return Map{
		"name":        p.Name,
		"value":       p.Value,
		"title":       p.Title,
		"type":        p.Type,
		"@type":       p.AtType,
		"unit":        p.Unit,
		"description": p.Description,
		"minimum":     p.Minimum,
		"maximum":     p.Maximum,
		"enum":        p.Enum,
		"readOnly":    p.ReadOnly,
		"multipleOf":  p.MultipleOf,
		"forms":       p.Forms,
		"deviceId":    p.DeviceId,
	}
}
