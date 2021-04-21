package addon

import (
	"fmt"
	"github.com/galenliu/gateway-addon/wot"
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
	*wot.PropertyAffordance

	//解决 PropertyAffordance中，未明确指定问题
	//Title       string `json:"title"`
	//AtType      string `json:"@type"`
	//Description string `json:"description,omitempty"`

	Name  string      `json:"name"`
	Value interface{} `json:"value,omitempty"`

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

	//err := json.UnmarshalFromString(description, &prop)
	//if err != nil {
	//	return nil
	//}
	//if prop.GetType() == TypeNumber {
	//	if gjson.Get(description, "minimum").Exists() || gjson.Get(description, "minimum").Exists() {
	//		schema := wot.NumberSchema{}
	//		if gjson.Get(description, "minimum").Exists() {
	//			schema.Minimum = gjson.Get(description, "minimum").Float()
	//		}
	//		if gjson.Get(description, "maximum").Exists() {
	//			schema.Minimum = gjson.Get(description, "maximum").Float()
	//		}
	//		if gjson.Get(description, "exclusiveMinimum").Exists() {
	//			schema.ExclusiveMinimum = gjson.Get(description, "exclusiveMinimum").Float()
	//		}
	//		if gjson.Get(description, "exclusiveMaximum").Exists() {
	//			schema.ExclusiveMaximum = gjson.Get(description, "exclusiveMaximum").Float()
	//		}
	//		if gjson.Get(description, "multipleOf").Exists() {
	//			schema.MultipleOf = gjson.Get(description, "multipleOf").Float()
	//		}
	//
	//	}
	//}
	//if prop.GetType() == TypeInteger {
	//	if gjson.Get(description, "minimum").Exists() || gjson.Get(description, "minimum").Exists() {
	//		schema := wot.IntegerSchema{}
	//		if gjson.Get(description, "minimum").Exists() {
	//			schema.Minimum = gjson.Get(description, "minimum").Int()
	//		}
	//		if gjson.Get(description, "maximum").Exists() {
	//			schema.Minimum = gjson.Get(description, "maximum").Int()
	//		}
	//		if gjson.Get(description, "exclusiveMinimum").Exists() {
	//			schema.ExclusiveMinimum = gjson.Get(description, "exclusiveMinimum").Int()
	//		}
	//		if gjson.Get(description, "exclusiveMaximum").Exists() {
	//			schema.ExclusiveMaximum = gjson.Get(description, "exclusiveMaximum").Int()
	//		}
	//		if gjson.Get(description, "multipleOf").Exists() {
	//			schema.MultipleOf = gjson.Get(description, "multipleOf").Int()
	//		}
	//
	//	}
	//}
	//if prop.GetType() == TypeString {
	//	if gjson.Get(description, "minLength").Exists() || gjson.Get(description, "maxLength").Exists() {
	//		schema := wot.StringSchema{}
	//		if gjson.Get(description, "minLength").Exists() {
	//			schema.MinLength = gjson.Get(description, "minLength").Int()
	//		}
	//		if gjson.Get(description, "maximum").Exists() {
	//			schema.MaxLength = gjson.Get(description, "maxLength").Int()
	//		}
	//
	//	}
	//}
	//
	e1 := json.UnmarshalFromString(description, &prop)
	if e1 != nil {
		return nil
	}
	typ := json.Get([]byte(description), "type").ToString()
	switch typ {
	case "array":
		var p wot.ArraySchema
		err := json.Unmarshal([]byte(description), &p)
		if err == nil {
			prop.IDataSchema = p
		}
	case "boolean":
		var p wot.BooleanSchema
		err := json.Unmarshal([]byte(description), &p)
		if err == nil {
			prop.IDataSchema = p
		}

	case "number":
		var p wot.NumberSchema
		err := json.Unmarshal([]byte(description), &p)
		if err == nil {
			prop.IDataSchema = p
		}

	case "integer":
		var p wot.IntegerSchema
		err := json.Unmarshal([]byte(description), &p)
		if err == nil {
			prop.IDataSchema = p
		}

	case "object":
		var p wot.ObjectSchema
		err := json.Unmarshal([]byte(description), &p)
		if err == nil {
			prop.IDataSchema = p
		}

	case "string":
		var p wot.StringSchema
		err := json.Unmarshal([]byte(description), &p)
		if err == nil {
			prop.IDataSchema = p
		}

	case "null":
		var p wot.NullSchema
		err := json.Unmarshal([]byte(description), &p)
		if err == nil {
			prop.IDataSchema = p
		}
	}

	prop.replyChan = make(chan Map)
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
	value = p.convert(value)
	switch p.GetType() {
	case TypeNumber:
		value = p.clampFloat(value.(float64))
	case TypeInteger:
		value = p.clampInt(int64(value.(int)))
	}
	if p.Value == value && !p.updateOnSameValue {
		return
	}
	if p.IsReadOnly() {
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
	switch p.GetType() {
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

func (p *Property) clampFloat(value float64) float64 {
	prop := p.IDataSchema.(*wot.NumberSchema)
	return prop.ClampFloat(value)
}

func (p *Property) clampInt(value int64) int64 {
	prop := p.IDataSchema.(*wot.IntegerSchema)
	return prop.ClampInt(value)
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

func (p *Property) AsDict() Map {
	return Map{
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
}
