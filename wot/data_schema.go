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
	Titles           []string      `json:"titles,omitempty"`
	Description      string        `json:"description,omitempty"`
	Descriptions     []string      `json:"descriptions,omitempty"`
	Unit             string        `json:"unit,omitempty"`
	Const            interface{}   `json:"const,omitempty"`
	OneOf            []DataSchema  `json:"oneOf,,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
	ReadOnly         bool          `json:"readOnly,omitempty"`
	WriteOnly        bool          `json:"writeOnly,omitempty"`
	Format           string        `json:"format,omitempty"`
	ContentEncoding  string        `json:"contentEncoding,,omitempty"`
	ContentMediaType string        `json:"contentMediaType,,omitempty"`

	Type string `json:"type"`
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
	Minimum          float64 `json:"minimum"`
	ExclusiveMinimum float64 `json:"exclusiveMinimum,omitempty"`
	Maximum          float64 `json:"maximum,omitempty"`
	ExclusiveMaximum float64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64 `json:"multipleOf,omitempty"`
}

func (n NumberSchema) ClampFloat(value float64) float64 {
	if n.Maximum != 0 {
		if value > n.Maximum {
			return n.Maximum
		}
	}
	if value < n.Minimum {
		return n.Minimum
	}
	return value
}

type IntegerSchema struct {
	*DataSchema
	Minimum          int64 `json:"minimum"`
	ExclusiveMinimum int64 `json:"exclusiveMinimum,omitempty"`
	Maximum          int64 `json:"maximum,omitempty"`
	ExclusiveMaximum int64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       int64 `json:"multipleOf,omitempty"`
}

func (n IntegerSchema) ClampInt(value int64) int64 {
	if n.Maximum != 0 {
		if value > n.Maximum {
			return n.Maximum
		}
	}
	if value < n.Minimum {
		return n.Minimum
	}
	return value
}

type ObjectSchema struct {
	*DataSchema
	Properties map[string]DataSchema `json:"properties,omitempty"`
	Required   string                `json:"required,omitempty"`
}

type StringSchema struct {
	*DataSchema
	MinLength int64 `json:"minLength,omitempty"`
	MaxLength int64 `json:"maxLength,omitempty"`
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

//func (d ArraySchema) MarshalJSON() ([]byte, error) {
//	return json.Marshal(d)
//}
//func (d BooleanSchema) MarshalJSON() ([]byte, error) {
//	return json.Marshal(d)
//}
//func (d NumberSchema) MarshalJSON() ([]byte, error) {
//	return json.Marshal(d)
//}
//func (d IntegerSchema) MarshalJSON() ([]byte, error) {
//	return json.Marshal(d)
//}
//func (d ObjectSchema) MarshalJSON() ([]byte, error) {
//	return json.Marshal(d)
//}
//func (d StringSchema) MarshalJSON() ([]byte, error) {
//	return json.Marshal(d)
//}
//func (d NullSchema) MarshalJSON() ([]byte, error) {
//	return json.Marshal(d)
//}

//func (d *ArraySchema) UnmarshalJSON(data []byte) error {
//	var this ArraySchema
//	var dataSchema DataSchema
//	err := json.Unmarshal(data, &dataSchema)
//	if err != nil {
//		return err
//	}
//	this.DataSchema = &dataSchema
//
//	if gjson.GetBytes(data, "items").Exists() {
//		var d []DataSchema
//		s := gjson.GetBytes(data, "items").String()
//		err := json.UnmarshalFromString(s, &d)
//		if err != nil {
//			this.Items = d
//		}
//	}
//
//	if gjson.GetBytes(data, "minItems").Exists() {
//		this.MinItems = int(gjson.GetBytes(data, "minItems").Int())
//	}
//
//	if gjson.GetBytes(data, "maxItems").Exists() {
//		this.MinItems = int(gjson.GetBytes(data, "maxItems").Int())
//	}
//
//	*d = this
//	return nil
//}
//func (d *BooleanSchema) UnmarshalJSON(data []byte) error {
//	var this BooleanSchema
//	var dataSchema DataSchema
//	err := json.Unmarshal(data, &dataSchema)
//	if err != nil {
//		return err
//	}
//	this.DataSchema = &dataSchema
//	*d = this
//	return nil
//}
//func (d *NumberSchema) UnmarshalJSON(data []byte) error {
//	var this NumberSchema
//	var dataSchema DataSchema
//	err := json.Unmarshal(data, &dataSchema)
//	if err != nil {
//		return err
//	}
//	this.DataSchema = &dataSchema
//
//	if gjson.GetBytes(data, "minimum").Exists() {
//		this.Minimum = float64(gjson.GetBytes(data, "minimum").Int())
//	}
//
//	if gjson.GetBytes(data, "exclusiveMinimum").Exists() {
//		this.ExclusiveMinimum = float64(gjson.GetBytes(data, "exclusiveMinimum").Int())
//	}
//	if gjson.GetBytes(data, "maximum").Exists() {
//		this.Maximum = float64(gjson.GetBytes(data, "maximum").Int())
//	}
//	if gjson.GetBytes(data, "exclusiveMaximum").Exists() {
//		this.ExclusiveMaximum = float64(gjson.GetBytes(data, "exclusiveMaximum").Int())
//	}
//	if gjson.GetBytes(data, "multipleOf").Exists() {
//		this.MultipleOf = float64(gjson.GetBytes(data, "multipleOf").Int())
//	}
//	*d = this
//	return nil
//}
//func (d *IntegerSchema) UnmarshalJSON(data []byte) error {
//	var this IntegerSchema
//	var dataSchema DataSchema
//	err := json.Unmarshal(data, &dataSchema)
//	if err != nil {
//		return err
//	}
//	this.DataSchema = &dataSchema
//
//	if gjson.GetBytes(data, "minimum").Exists() {
//		this.Minimum = int64(gjson.GetBytes(data, "minimum").Int())
//	}
//
//	if gjson.GetBytes(data, "exclusiveMinimum").Exists() {
//		this.ExclusiveMinimum = int64(gjson.GetBytes(data, "exclusiveMinimum").Int())
//	}
//	if gjson.GetBytes(data, "maximum").Exists() {
//		this.Maximum = int64(gjson.GetBytes(data, "maximum").Int())
//	}
//	if gjson.GetBytes(data, "exclusiveMaximum").Exists() {
//		this.ExclusiveMaximum = int64(gjson.GetBytes(data, "exclusiveMaximum").Int())
//	}
//	if gjson.GetBytes(data, "multipleOf").Exists() {
//		this.MultipleOf = int64(gjson.GetBytes(data, "multipleOf").Int())
//	}
//
//	*d = this
//	return nil
//}
//func (d *ObjectSchema) UnmarshalJSON(data []byte) error {
//	var this ObjectSchema
//	var dataSchema DataSchema
//	err := json.Unmarshal(data, &dataSchema)
//	if err != nil {
//		return err
//	}
//	this.DataSchema = &dataSchema
//
//	if gjson.GetBytes(data, "required").Exists() {
//		this.Required = gjson.GetBytes(data, "required").String()
//	}
//
//	if gjson.GetBytes(data, "properties").Exists() {
//		var d map[string]DataSchema
//		s := gjson.GetBytes(data, "items").String()
//		err := json.UnmarshalFromString(s, &d)
//		if err != nil {
//			this.Properties = d
//		}
//	}
//	*d = this
//	return nil
//}
//func (d *StringSchema) UnmarshalJSON(data []byte) error {
//	var this StringSchema
//	var dataSchema DataSchema
//	err := json.Unmarshal(data, &dataSchema)
//	if err != nil {
//		return err
//	}
//	this.DataSchema = &dataSchema
//
//	if gjson.GetBytes(data, "minLength").Exists() {
//		this.MinLength = int64(gjson.GetBytes(data, "minLength").Int())
//	}
//	if gjson.GetBytes(data, "maxLength").Exists() {
//		this.MaxLength = int64(gjson.GetBytes(data, "maxLength").Int())
//	}
//	*d = this
//	return nil
//}
//func (d *NullSchema) UnmarshalJSON(data []byte) error {
//	var this NullSchema
//	var dataSchema DataSchema
//	err := json.Unmarshal(data, &dataSchema)
//	if err != nil {
//		return err
//	}
//	this.DataSchema = &dataSchema
//	*d = this
//	return nil
//}

type IDataSchema interface {
	GetType() string
	SetAtType(string)

	IsReadOnly() bool
	SetType(string)
	SetTitle(s string)
	//MarshalJSON() ([]byte, error)
	//UnmarshalJSON(data []byte) error
}
