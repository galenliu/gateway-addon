package wot

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/xiam/to"
)

type PropertyAffordance struct {
	*InteractionAffordance
	IDataSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value,omitempty"`
}

func NewPropertyAffordanceFromString(description string) *PropertyAffordance {
	var p = PropertyAffordance{}
	p.InteractionAffordance = NewInteractionAffordanceFromString(description)
	p.IDataSchema = NewDataSchemaFromString(description)
	if gjson.Get(description, "observable").Exists() {
		p.Observable = gjson.Get(description, "observable").Bool()
	}
	return &p
}

func (p *PropertyAffordance) SetCachedValue(value interface{}) {
	value = p.convert(value)
	p.Value = p.clamp(value)
}

//把属性值转换成确定的类型
func (p *PropertyAffordance) convert(v interface{}) interface{} {
	switch p.GetType() {
	case Number:
		return to.Float64(v)
	case Integer:
		return int(to.Uint64(v))
	case Boolean:
		return to.Bool(v)
	default:
		return v
	}
}

//保证属性值在允许的范围
func (p *PropertyAffordance) clamp(v interface{}) interface{} {
	switch p.IDataSchema.(type) {
	case *NumberSchema:
		dd := p.IDataSchema.(*NumberSchema)
		return dd.ClampFloat(to.Float64(v))
	case NumberSchema:
		dd := p.IDataSchema.(*NumberSchema)
		return dd.ClampFloat(to.Float64(v))
	case *IntegerSchema:
		dd := p.IDataSchema.(*IntegerSchema)
		return dd.ClampInt(to.Int64(v))
	case IntegerSchema:
		dd := p.IDataSchema.(*IntegerSchema)
		return dd.ClampInt(to.Int64(v))
	default:
		return v
	}

}

func (p PropertyAffordance) MarshalJSON() ([]byte, error) {

	//log.Println("property Data Schema type:%s", reflect.TypeOf(p.IDataSchema))
	if p.IDataSchema == nil {
		return nil, fmt.Errorf("dataschema err")
	}
	switch p.IDataSchema.(type) {
	case *ArraySchema:
		dd := p.IDataSchema.(*ArraySchema)
		var pa = ArrayPropertyAffordance{p.InteractionAffordance, dd, p.Observable, p.Value}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case ArraySchema:
		dd := p.IDataSchema.(ArraySchema)
		var pa = ArrayPropertyAffordance{p.InteractionAffordance, &dd, p.Observable, p.Value}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *BooleanSchema:
		dd := p.IDataSchema.(*BooleanSchema)
		var pa = BooleanPropertyAffordance{p.InteractionAffordance, dd, p.Observable, to.Bool(p.Value)}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case BooleanSchema:
		dd := p.IDataSchema.(BooleanSchema)
		var pa = BooleanPropertyAffordance{p.InteractionAffordance, &dd, p.Observable, to.Bool(p.Value)}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *NumberSchema:
		dd := p.IDataSchema.(*NumberSchema)
		var pa = NumberPropertyAffordance{p.InteractionAffordance, dd, p.Observable, to.Float64(p.Value)}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case NumberSchema:
		dd := p.IDataSchema.(NumberSchema)
		var pa = NumberPropertyAffordance{p.InteractionAffordance, &dd, p.Observable, to.Float64(p.Value)}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *IntegerSchema:
		dd := p.IDataSchema.(*IntegerSchema)
		var pa = IntegerPropertyAffordance{p.InteractionAffordance, dd, p.Observable, to.Int64(p.Value)}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case IntegerSchema:
		dd := p.IDataSchema.(IntegerSchema)
		var pa = IntegerPropertyAffordance{p.InteractionAffordance, &dd, p.Observable, to.Int64(p.Value)}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *ObjectSchema:
		dd := p.IDataSchema.(*ObjectSchema)
		var pa = ObjectPropertyAffordance{p.InteractionAffordance, dd, p.Observable, p.Value}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case ObjectSchema:
		dd := p.IDataSchema.(ObjectSchema)
		var pa = ObjectPropertyAffordance{p.InteractionAffordance, &dd, p.Observable, p.Value}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *StringSchema:
		dd := p.IDataSchema.(*StringSchema)
		var pa = StringPropertyAffordance{p.InteractionAffordance, dd, p.Observable, to.String(p.Value)}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case StringSchema:
		dd := p.IDataSchema.(StringSchema)
		var pa = StringPropertyAffordance{p.InteractionAffordance, &dd, p.Observable, to.String(p.Value)}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *NullSchema:
		dd := p.IDataSchema.(*NullSchema)
		var pa = NullPropertyAffordance{p.InteractionAffordance, dd, p.Observable, p.Value}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case NullSchema:
		dd := p.IDataSchema.(NullSchema)
		var pa = NullPropertyAffordance{p.InteractionAffordance, &dd, p.Observable, p.Value}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	default:
		return nil, fmt.Errorf("property type err")
	}
}

type ArrayPropertyAffordance struct {
	*InteractionAffordance
	*ArraySchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

type BooleanPropertyAffordance struct {
	*InteractionAffordance
	*BooleanSchema
	Observable bool `json:"observable"`
	Value      bool `json:"value"`
}

type NumberPropertyAffordance struct {
	*InteractionAffordance
	*NumberSchema
	Observable bool    `json:"observable"`
	Value      float64 `json:"value"`
}

type IntegerPropertyAffordance struct {
	*InteractionAffordance
	*IntegerSchema
	Observable bool  `json:"observable"`
	Value      int64 `json:"value"`
}

type ObjectPropertyAffordance struct {
	*InteractionAffordance
	*ObjectSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

type StringPropertyAffordance struct {
	*InteractionAffordance
	*StringSchema
	Observable bool   `json:"observable"`
	Value      string `json:"value"`
}

type NullPropertyAffordance struct {
	*InteractionAffordance
	*NullSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

func (p *PropertyAffordance) SetMinValue(v interface{}) {
	switch p.IDataSchema.(type) {
	case *NumberSchema:
		dd := p.IDataSchema.(*NumberSchema)
		dd.Minimum = to.Float64(v)
	case NumberSchema:
		dd := p.IDataSchema.(*NumberSchema)
		dd.Minimum = to.Float64(v)
	case *IntegerSchema:
		dd := p.IDataSchema.(*IntegerSchema)
		dd.Minimum = to.Int64(v)
	case IntegerSchema:
		dd := p.IDataSchema.(*IntegerSchema)
		dd.Minimum = to.Int64(v)

	default:
		fmt.Print("property type err")
		return
	}
}

func (p *PropertyAffordance) GetDefaultValue() interface{} {
	switch p.IDataSchema.(type) {
	case *NumberSchema:
		dd := p.IDataSchema.(*NumberSchema)
		return dd.Minimum
	case NumberSchema:
		dd := p.IDataSchema.(NumberSchema)
		return dd.Minimum
	case *IntegerSchema:
		dd := p.IDataSchema.(*IntegerSchema)
		return dd.Minimum
	case IntegerSchema:
		dd := p.IDataSchema.(IntegerSchema)
		return dd.Minimum

	case *BooleanSchema:
		return false
	case BooleanSchema:
		return false

	case *StringSchema:
		return ""
	case StringSchema:
		return ""

	default:
		return nil
	}
}

func (p *PropertyAffordance) SetMaxValue(v interface{}) {
	switch p.IDataSchema.(type) {
	case *NumberSchema:
		dd := p.IDataSchema.(*NumberSchema)
		dd.Maximum = to.Float64(v)
	case NumberSchema:
		dd := p.IDataSchema.(*NumberSchema)
		dd.Maximum = to.Float64(v)
	case *IntegerSchema:
		dd := p.IDataSchema.(*IntegerSchema)
		dd.Maximum = to.Int64(v)
	case IntegerSchema:
		dd := p.IDataSchema.(*IntegerSchema)
		dd.Maximum = to.Int64(v)
	default:
		fmt.Print("property type err")
		return
	}
}
