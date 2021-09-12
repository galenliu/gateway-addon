package properties

import json "github.com/json-iterator/go"

type Property struct {
	Name        string        `json:"name"`
	Title       string        `json:"title"`
	Type        string        `json:"type"`
	AtType      string        `json:"@type"`
	Unit        string        `json:"unit"`
	Description string        `json:"description"`
	Minimum     interface{}   `json:"minimum"`
	Maximum     interface{}   `json:"maximum"`
	Enum        []interface{} `json:"enum"`
	ReadOnly    bool          `json:"readOnly"`
	MultipleOf  interface{}   `json:"multipleOf"`
	Forms       []interface{} `json:"forms"`
	Value       interface{}   `json:"value"`
	DeviceId    string        `json:"deviceId"`
}

func NerPropertyFormString(s string) *Property {
	var p Property
	_ = json.UnmarshalFromString(s, &p)
	return &p

}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) GetTitle() string {
	return p.Title
}

func (p *Property) GetType() string {
	return p.Type
}

func (p *Property) GetAtType() string {
	return p.AtType
}

func (p *Property) GetUnit() string {
	return p.Unit
}

func (p *Property) GetDescription() string {
	return p.Description
}

func (p *Property) GetMinimum() interface{} {
	return p.Minimum
}
func (p *Property) GetMaximum() interface{} {
	return p.Maximum
}

func (p *Property) IsReadOnly() bool {
	return p.ReadOnly
}

func (p *Property) GetMultipleOf() interface{} {
	return p.MultipleOf
}

func (p *Property) GetForms() []interface{} {
	return p.Forms
}

func (p *Property) GetValue() interface{} {
	return p.Value
}
