package properties

type NumberProperty struct {
	*addon.Property
}

func NewNumberProperty(typ string) *NumberProperty {
	p := addon.NewProperty(typ)
	p.SetType(TypeNumber)
	return &NumberProperty{p}
}

// SetValue sets a value
func (prop *NumberProperty) SetValue(value float64) {
	prop.UpdateValue(value)
}

// GetValue returns the value as bool
func (prop *NumberProperty) GetValue() float64 {
	return prop.Property.GetValue().(float64)
}

// OnValueRemoteGet calls fn when the value was read by a client.
func (prop *NumberProperty) OnValueRemoteGet(fn func() float64) {
	prop.OnValueGet(func() interface{} {
		return fn()
	})
}

// OnValueRemoteUpdate calls fn when the value was updated by a client.
func (prop *NumberProperty) OnValueRemoteUpdate(fn func(float64)) {
	prop.OnValueUpdate(func(property *addon.Property, newValue, oldValue interface{}) {
		fn(newValue.(float64))
	})
}
