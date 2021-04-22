package properties

type IntegerProperty struct {
	*addon.Property
}

func NewIntegerProperty(typ string) *IntegerProperty {
	p := addon.NewProperty(typ)
	p.SetType(TypeInteger)
	return &IntegerProperty{p}
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

func (prop *IntegerProperty) SetMinValue(v int64) {
	prop.Property.SetMinValue(v)
}

func (prop *IntegerProperty) SetMaxValue(v int64) {
	prop.Property.SetMaxValue(v)
}
