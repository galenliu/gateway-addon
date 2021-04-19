package properties

const TypeColorTemperatureProperty = "ColorTemperatureProperty"

type ColorTemperatureProperty struct {
	*IntegerProperty
}

func NewColorTemperatureProperty() *ColorTemperatureProperty {
	p := NewIntegerProperty(TypeColorTemperatureProperty)
	p.SetType(TypeInteger)
	p.Name = ColorTemperature
	p.SetValue(0)
	p.Unit = UnitKelvin

	return &ColorTemperatureProperty{p}
}
