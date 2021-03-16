package addon

import (
	"github.com/xiam/to"
	//json "github.com/json-iterator/go"
)

type ChangeFunc func(property *Property, newValue, oldValue interface{})
type GetFunc func() interface{}

type Property struct {
	AtType      string      `json:"@type"` //引用的类型
	Type        string      `json:"type"`  //数据的格式
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Name        string      `json:"name"`
	ReadOnly    bool        `json:"readOnly"`
	Visible     bool        `json:"visible"`
	Value       interface{} `json:"value"`

	Unit      string      `json:"unit,omitempty"`
	Minimum   interface{} `json:"minimum,omitempty"`
	Maximum   interface{} `json:"maximum,omitempty"`
	StepValue interface{} `json:"stepValue,omitempty"`

	Enum []string `json:"enum,omitempty"`

	DeviceId string `json:"-"`

	updateOnSameValue bool

	valueChangeFuncs []ChangeFunc
	valueGetFunc     GetFunc
}

func NewProperty(typ string) *Property {
	prop := &Property{
		AtType:           typ,
		valueChangeFuncs: make([]ChangeFunc, 0),
	}
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

func (prop *Property) Update(new *Property) {
	prop.Value = new.Value
	prop.Name = new.Name
	prop.Type = new.Type
	prop.ReadOnly = new.ReadOnly
	prop.Title = new.Title
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
