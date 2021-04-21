package properties

import "github.com/galenliu/gateway-addon"

type StringProperty struct {
	*addon.Property
}

func NewStringProperty(typ string) *StringProperty {
	p := addon.NewProperty(typ)
	p.SetType(TypeString)
	return &StringProperty{p}
}

// SetValue sets a value
func (prop *StringProperty) SetValue(value string) {
	prop.UpdateValue(value)
}

// GetValue returns the value as bool
func (prop *StringProperty) GetValue() string {
	return prop.Property.GetValue().(string)
}

// OnValueRemoteGet calls fn when the value was read by a client.
func (prop *StringProperty) OnValueRemoteGet(fn func() string) {
	prop.OnValueGet(func() interface{} {
		return fn()
	})
}

// OnValueRemoteUpdate calls fn when the value was updated by a client.
func (prop *StringProperty) OnValueRemoteUpdate(fn func(string)) {
	prop.OnValueUpdate(func(property *addon.Property, newValue, oldValue interface{}) {
		fn(newValue.(string))
	})
}
