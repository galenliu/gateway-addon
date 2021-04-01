package addon

import (
	"fmt"
	json "github.com/json-iterator/go"
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
	AtType      string      `json:"@type"` //引用的类型(OnOffProperty ...)
	Type        string      `json:"type"`  //数据的格式(string,boolean ...)
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Name        string      `json:"name"`
	ReadOnly    bool        `json:"readOnly"`
	Visible     bool        `json:"visible"`
	Value       interface{} `json:"value"`

	Unit       string      `json:"unit,omitempty"`
	Minimum    interface{} `json:"minimum,omitempty"`
	Maximum    interface{} `json:"maximum,omitempty"`
	MultipleOf int         `json:"multipleOf,omitempty"`

	StepValue interface{} `json:"stepValue,omitempty"`

	Enum []interface{} `json:"enum,omitempty"`

	DeviceId string `json:"deviceId"`

	device Owner

	updateOnSameValue bool

	valueChangeFuncs []ChangeFunc
	valueGetFunc     GetFunc

	verbose bool
}

func NewProperty(typ string) *Property {
	prop := &Property{
		AtType:           typ,
		valueChangeFuncs: make([]ChangeFunc, 0),
	}
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
