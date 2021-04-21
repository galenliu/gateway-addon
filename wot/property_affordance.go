package wot

import (
	"fmt"
	json "github.com/json-iterator/go"
)

type PropertyAffordance struct {
	*InteractionAffordance
	IDataSchema
	Observable bool `json:"observable"`
}

func (p PropertyAffordance) MarshalJSON() ([]byte, error) {

	if p.IDataSchema == nil {
		return nil, fmt.Errorf("dataschema err")
	}
	switch p.IDataSchema.(type) {
	case *ArraySchema:
		dd := p.IDataSchema.(*ArraySchema)
		var pa = ArrayPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *BooleanSchema:
		dd := p.IDataSchema.(*BooleanSchema)
		var pa = BooleanPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *NumberSchema:
		dd := p.IDataSchema.(*NumberSchema)
		var pa = NumberPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *IntegerSchema:
		dd := p.IDataSchema.(*IntegerSchema)
		var pa = IntegerPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *ObjectSchema:
		dd := p.IDataSchema.(*ObjectSchema)
		var pa = ObjectPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *StringSchema:
		dd := p.IDataSchema.(*StringSchema)
		var pa = StringPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *NullSchema:
		dd := p.IDataSchema.(*NullSchema)
		var pa = NullPropertyAffordance{p.InteractionAffordance, dd, p.Observable}
		pa.AtType = dd.AtType
		return json.MarshalIndent(pa, "", "  ")
	default:
		return nil, fmt.Errorf("property type err")
	}
}

//func (p *PropertyAffordance) UnmarshalJSON(data []byte) error {
//
//	typ := json.Get(data, "type").ToString()
//
//	switch typ {
//	case "array":
//		var prop ArrayPropertyAffordance
//		err := json.Unmarshal(data, &prop)
//		if err != nil {
//			return err
//		}
//		*p = PropertyAffordance{
//			InteractionAffordance: prop.InteractionAffordance,
//			IDataSchema:           prop.ArraySchema,
//			Observable:            prop.Observable,
//		}
//
//	case "boolean":
//		var prop BooleanPropertyAffordance
//		err := json.Unmarshal(data, &prop)
//		if err != nil {
//			return err
//		}
//		p = &PropertyAffordance{
//			InteractionAffordance: prop.InteractionAffordance,
//			IDataSchema:           prop.BooleanSchema,
//			Observable:            prop.Observable,
//		}
//
//	case "number":
//		var prop NumberPropertyAffordance
//		err := json.Unmarshal(data, &prop)
//		if err != nil {
//			return err
//		}
//		p = &PropertyAffordance{
//			InteractionAffordance: prop.InteractionAffordance,
//			IDataSchema:           prop.NumberSchema,
//			Observable:            prop.Observable,
//		}
//
//	case "integer":
//		var prop IntegerPropertyAffordance
//		err := json.Unmarshal(data, &prop)
//		if err != nil {
//			return err
//		}
//		p = &PropertyAffordance{
//			InteractionAffordance: prop.InteractionAffordance,
//			IDataSchema:           prop.IntegerSchema,
//			Observable:            prop.Observable,
//		}
//
//	case "object":
//		var prop ObjectPropertyAffordance
//		err := json.Unmarshal(data, &prop)
//		if err != nil {
//			return err
//		}
//		p = &PropertyAffordance{
//			InteractionAffordance: prop.InteractionAffordance,
//			IDataSchema:           prop.ObjectSchema,
//			Observable:            prop.Observable,
//		}
//
//	case "string":
//		var prop StringPropertyAffordance
//		err := json.Unmarshal(data, &prop)
//		if err != nil {
//			return err
//		}
//		p = &PropertyAffordance{
//			InteractionAffordance: prop.InteractionAffordance,
//			IDataSchema:           prop.StringSchema,
//			Observable:            prop.Observable,
//		}
//
//	case "null":
//		var prop NullPropertyAffordance
//		err := json.Unmarshal(data, &prop)
//		if err != nil {
//			return err
//		}
//		p = &PropertyAffordance{
//			InteractionAffordance: prop.InteractionAffordance,
//			IDataSchema:           prop.NullSchema,
//			Observable:            prop.Observable,
//		}
//	}
//	return nil
//}

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
