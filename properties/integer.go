package properties

import "addon"

type IntegerProperty struct {
	*addon.Property
}

func NewIntegerProperty(typ string) *IntegerProperty {
	integer := addon.NewProperty(typ)
	integer.Type = TypeInteger
	return &IntegerProperty{integer}
}

// SetValue sets a value
func (prop *IntegerProperty) SetCachedValueAndNotify(value int) {
	prop.Property.SetCachedValueAndNotify(value)
}

// GetValue returns the value as bool
func (prop *IntegerProperty) GetValue() int {
	return prop.Property.GetValue().(int)
}

// OnValueRemoteGet calls fn when the value was read by a client.
func (prop *IntegerProperty) OnValueRemoteGet(fn func() int) {

}

// OnValueRemoteUpdate calls fn when the value was updated by a client.
func (prop *IntegerProperty) OnValueRemoteUpdate(fn func(int)) {
	prop.OnValueUpdate(func(property *addon.Property, newValue, oldValue interface{}) {
		fn(newValue.(int))
	})
}

func (prop *IntegerProperty) SetMinValue(value int) {
	prop.Minimum = value
}

func (prop *IntegerProperty) SetMaxValue(value int) {
	prop.Maximum = value
}

func (prop *IntegerProperty) SetStepValue(value int) {
	prop.StepValue = value
}
