package wot

const (
	JSON      = "application/json"
	LdJSON    = "application/ld+json"
	SenmlJSON = "application/senml+json"
	CBOR      = "application/cbor"
	SenmlCbor = "application/senml+cbor"

	XML      = "application/xml"
	SenmlXML = "application/senml+xml"
	EXI      = "application/exi"
)

type DataSchema struct {
	*InteractionAffordance
	Type  string      `json:"type"`
	Const interface{} `json:"const,omitempty"`
	Unit  string      `json:"unit,omitempty"`

	OneOf []DataSchema  `json:"oneOf,,omitempty"`
	Enum  []interface{} `json:"enum,omitempty"`

	ReadOnly  bool `json:"readOnly"`
	WriteOnly bool `json:"writeOnly"`

	Format           string `json:"format,omitempty"`
	ContentEncoding  string `json:"contentEncoding,,omitempty"`
	ContentMediaType string `json:"contentMediaType,,omitempty"`
}

type ArraySchema struct {
	Items    []DataSchema `json:"items,omitempty"`
	MinItems int          `json:"minItems,omitempty"`
	maxItems int          `json:"maxItems,omitempty"`
}

type NumberSchema struct {
	Minimum          float64 `json:"minimum,omitempty"`
	ExclusiveMinimum float64 `json:"exclusiveMinimum,omitempty"`
	Maximum          float64 `json:"maximum,omitempty"`
	ExclusiveMaximum float64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64 `json:"multipleOf,omitempty"`
}

type IntegerSchema struct {
	Minimum          int64 `json:"minimum,omitempty"`
	ExclusiveMinimum int64 `json:"exclusiveMinimum,omitempty"`
	Maximum          int64 `json:"maximum,omitempty"`
	ExclusiveMaximum int64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       int64 `json:"multipleOf,omitempty"`
}

type ObjectSchema struct {
	Properties map[string]DataSchema `json:"properties"`
	Required   []string              `json:"required"`
}

type StringSchema struct {
	MinLength int64 `json:"minLength"`
	MaxLength int64 `json:"maxLength"`
}
