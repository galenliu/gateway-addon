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
	AtType           string        `json:"@type"`
	Title            string        `json:"title"`
	Titles           []string      `json:"titles"`
	Description      string        `json:"description"`
	Descriptions     []string      `json:"descriptions"`
	Unit             string        `json:"unit,omitempty"`
	Const            interface{}   `json:"const,omitempty"`
	OneOf            []DataSchema  `json:"oneOf,,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
	ReadOnly         bool          `json:"readOnly"`
	WriteOnly        bool          `json:"writeOnly"`
	Format           string        `json:"format,omitempty"`
	ContentEncoding  string        `json:"contentEncoding,,omitempty"`
	ContentMediaType string        `json:"contentMediaType,,omitempty"`

	Type   string      `json:"type"`
	Schema interface{} `json:"schema"`
}

type ArraySchema struct {
	*DataSchema
	Items    []DataSchema `json:"items,omitempty"`
	MinItems int          `json:"minItems,omitempty"`
	MaxItems int          `json:"maxItems,omitempty"`
}

type BooleanSchema struct {
	*DataSchema
}

type NumberSchema struct {
	*DataSchema
	Minimum          float64 `json:"minimum,omitempty"`
	ExclusiveMinimum float64 `json:"exclusiveMinimum,omitempty"`
	Maximum          float64 `json:"maximum,omitempty"`
	ExclusiveMaximum float64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64 `json:"multipleOf,omitempty"`
}

type IntegerSchema struct {
	*DataSchema
	Minimum          int64 `json:"minimum,omitempty"`
	ExclusiveMinimum int64 `json:"exclusiveMinimum,omitempty"`
	Maximum          int64 `json:"maximum,omitempty"`
	ExclusiveMaximum int64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       int64 `json:"multipleOf,omitempty"`
}

type ObjectSchema struct {
	*DataSchema
	Properties map[string]DataSchema `json:"properties"`
	Required   string                `json:"required"`
}

type StringSchema struct {
	*DataSchema
	MinLength int64 `json:"minLength"`
	MaxLength int64 `json:"maxLength"`
}

type NullSchema struct {
	*DataSchema
}

func (d *DataSchema) GetType() string {
	return d.Type
}

func (d *DataSchema) SetAtType(s string) {
	d.AtType = s
}

func (d *DataSchema) SetType(s string) {
	d.Type = s
}

func (d *DataSchema) SetTitle(s string) {
	d.Title = s
}

func (d *DataSchema) IsReadOnly() bool {
	return d.ReadOnly
}

type IDataSchema interface {
	GetType() string
	SetAtType(string)
	IsReadOnly() bool
	SetType(string)
	SetTitle(s string)
}
