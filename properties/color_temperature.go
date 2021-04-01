package properties

const TypeColorTemperatureProperty = "ColorTemperatureProperty"

type ColorTemperatureProperty struct {
	*IntegerProperty
}

func NewColorTemperatureProperty() *ColorTemperatureProperty {
	p := NewIntegerProperty(TypeColorTemperatureProperty)
	p.Type = TypeInteger
	p.Name = ColorTemperature
	p.SetStepValue(1)
	p.SetValue(0)
	p.Unit = UnitKelvin

	return &ColorTemperatureProperty{p}
}
