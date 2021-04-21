package wot

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type PropertyAffordance struct {
	*InteractionAffordance
	IDataSchema
	Observable bool `json:"observable"`
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

func (p PropertyAffordance) MarshalJSON() ([]byte, error) {

	//log.Println("property Data Schema type:%s", reflect.TypeOf(p.IDataSchema))
	if p.IDataSchema == nil {
		return nil, fmt.Errorf("dataschema err")
	}
	switch p.IDataSchema.(type) {
	case *ArraySchema:
		dd := p.IDataSchema.(*ArraySchema)
		var pa = ArrayPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case ArraySchema:
		dd := p.IDataSchema.(ArraySchema)
		var pa = ArrayPropertyAffordance{p.InteractionAffordance, &dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *BooleanSchema:
		dd := p.IDataSchema.(*BooleanSchema)
		var pa = BooleanPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case BooleanSchema:
		dd := p.IDataSchema.(BooleanSchema)
		var pa = BooleanPropertyAffordance{p.InteractionAffordance, &dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *NumberSchema:
		dd := p.IDataSchema.(*NumberSchema)
		var pa = NumberPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case NumberSchema:
		dd := p.IDataSchema.(NumberSchema)
		var pa = NumberPropertyAffordance{p.InteractionAffordance, &dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *IntegerSchema:
		dd := p.IDataSchema.(*IntegerSchema)
		var pa = IntegerPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case IntegerSchema:
		dd := p.IDataSchema.(IntegerSchema)
		var pa = IntegerPropertyAffordance{p.InteractionAffordance, &dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *ObjectSchema:
		dd := p.IDataSchema.(*ObjectSchema)
		var pa = ObjectPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case ObjectSchema:
		dd := p.IDataSchema.(ObjectSchema)
		var pa = ObjectPropertyAffordance{p.InteractionAffordance, &dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *StringSchema:
		dd := p.IDataSchema.(*StringSchema)
		var pa = StringPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case StringSchema:
		dd := p.IDataSchema.(StringSchema)
		var pa = StringPropertyAffordance{p.InteractionAffordance, &dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *NullSchema:
		dd := p.IDataSchema.(*NullSchema)
		var pa = NullPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case NullSchema:
		dd := p.IDataSchema.(NullSchema)
		var pa = NullPropertyAffordance{p.InteractionAffordance, &dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	default:
		return nil, fmt.Errorf("property type err")
	}
}

type ArrayPropertyAffordance struct {
	*InteractionAffordance
	*ArraySchema
	Observable bool `json:"observable"`
}

type BooleanPropertyAffordance struct {
	*InteractionAffordance
	*BooleanSchema
	Observable bool `json:"observable"`
}

type NumberPropertyAffordance struct {
	*InteractionAffordance
	*NumberSchema
	Observable bool `json:"observable"`
}

type IntegerPropertyAffordance struct {
	*InteractionAffordance
	*IntegerSchema
	Observable bool `json:"observable"`
}

type ObjectPropertyAffordance struct {
	*InteractionAffordance
	*ObjectSchema
	Observable bool `json:"observable"`
}

type StringPropertyAffordance struct {
	*InteractionAffordance
	*StringSchema
	Observable bool `json:"observable"`
}

type NullPropertyAffordance struct {
	*InteractionAffordance
	*NullSchema
	Observable bool `json:"observable"`
}
