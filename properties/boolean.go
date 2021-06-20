package properties

import (
	addon "github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/wot/definitions/data_schema"
)

type BooleanProperty struct {
	*addon.Property
}

func NewBooleanProperty(typ string) *BooleanProperty {
	boolean := &addon.Property{}
	boolean.IDataSchema = data_schema.NewBooleanSchema()
	boolean.SetAtType(typ)
	return &BooleanProperty{boolean}
}

// SetValue sets a value
func (prop *BooleanProperty) SetBoolenValue(value bool) {
	prop.UpdateValue(value)
}

// GetValue returns the value as bool
func (prop *BooleanProperty) GetValue() bool {
	return prop.Property.GetValue().(bool)
}

// OnValueRemoteGet calls fn when the value was read by a client.
func (prop *BooleanProperty) OnValueRemoteGet(fn func() bool) {
	prop.OnValueGet(func() interface{} {
		return fn()
	})
}

// OnValueRemoteUpdate calls fn when the value was updated by a client.
func (prop *BooleanProperty) OnValueRemoteUpdate(fn func(bool)) {
	prop.OnValueUpdate(func(property *addon.Property, newValue, oldValue interface{}) {
		fn(newValue.(bool))
	})
}
