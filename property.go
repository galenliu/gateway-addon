package addon

import (
	"addon/wot"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/xiam/to"
	"log"
	//json "github.com/json-iterator/go"
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

	Schema interface{}

	DeviceId string `json:"deviceId,omitempty"`

	device            Owner
	updateOnSameValue bool

	valueChangeFuncs []ChangeFunc
	valueGetFunc     GetFunc
	verbose          bool
}

func NewPropertyFromString(description string) *Property {
	var prop Property
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

func (prop *Property) OnValueUpdate(fn ChangeFunc) {
	prop.valueChangeFuncs = append(prop.valueChangeFuncs, fn)
}

func (prop *Property) OnValueGet(fn GetFunc) {
	prop.valueGetFunc = fn
}

func (prop *Property) GetValue() interface{} {
	return prop.getValue()
}

func (prop *Property) getValue() interface{} {
	if prop.valueGetFunc != nil {
		prop.UpdateValue(prop.valueGetFunc())
	}
	return prop.Value
}

func (prop *Property) SetCachedValueAndNotify(value interface{}) {
	prop.UpdateValue(value)
	data := make(map[string]interface{})
	data["property"] = prop.ToString()
	prop.device.Send(DevicePropertyChangedNotification, data)
}

func (prop *Property) UpdateValue(value interface{}) {
	value = prop.convert(value)
	switch prop.Type {
	case TypeNumber:
		value = prop.clampFloat(value.(float64))
	case TypeInteger:
		value = prop.clampInt(value.(int))
	}
	if prop.Value == value && !prop.updateOnSameValue {
		return
	}
	if prop.ReadOnly {
		return
	}
	oldValue := prop.Value
	prop.Value = value
	prop.onValueUpdate(prop.valueChangeFuncs, value, oldValue)
}

func (prop *Property) onValueUpdate(funcs []ChangeFunc, newValue, oldValue interface{}) {
	for _, fn := range funcs {
		fn(prop, newValue, oldValue)
	}

}

func (prop *Property) Update(d []byte) {
	title := json.Get(d, "title").ToString()
	if title != "" && prop.Title != title {
		prop.Title = title
	}
	value := json.Get(d, "value").GetInterface()
	if value != nil && prop.Value != value {
		prop.Value = value
	}
}

func (prop *Property) convert(v interface{}) interface{} {
	switch prop.Type {
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

func (prop *Property) clampFloat(value float64) interface{} {
	min, minOK := prop.Minimum.(float64)
	max, maxOK := prop.Maximum.(float64)
	if maxOK == true && value > max {
		value = max
	} else if minOK == true && value < min {
		value = min
	}
	return value
}

func (prop *Property) clampInt(value int) interface{} {
	min, minOK := prop.Minimum.(int)
	max, maxOK := prop.Maximum.(int)
	if maxOK == true && value > max {
		value = max
	} else if minOK == true && value < min {
		value = min
	}
	return value
}

func (prop *Property) SetValue(newValue interface{}) {
	if prop.verbose {
		log.Printf("property(%s) set value not imp", prop.GetName())
	}
}

func (prop *Property) GetName() string {
	return prop.Name
}

func (prop *Property) SetName(name string) {
	prop.Name = name
}

func (prop *Property) GetAtType() string {
	return prop.AtType
}

func (prop *Property) GetType() string {
	return prop.Type
}

//func (prop *Property) MarshalJSON() ([]byte, error) {
//	return json.MarshalIndent(prop, "", " ")
//}

func (prop *Property) ToString() string {
	str, err := json.MarshalToString(prop)
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}
	return str
}

func (prop *Property) SetOwner(owner Owner) {
	prop.device = owner
}

type NotifyProp struct {
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
	Title    string      `json:"title"`
	Type     string      `json:"type"`
	AtType   string      `json:"@type"`
	DeviceId string      `json:"deviceId"`
}

func (prop *Property) GetNotifyDescription() []byte {
	p := &NotifyProp{}
	p.Name = prop.Name
	p.Value = prop.Value
	p.DeviceId = prop.DeviceId
	p.Title = prop.Title
	p.Type = prop.Type
	p.AtType = prop.AtType
	d, _ := json.Marshal(p)
	return d
}
