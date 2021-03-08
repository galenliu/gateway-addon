package properties

const TypeBrightnessProperty = "BrightnessProperty"

type BrightnessProperty struct {
	*IntegerProperty
}

func NewBrightnessProperty() *BrightnessProperty {
	p := NewIntegerProperty(TypeBrightnessProperty)
	p.Type = TypeInteger
	p.Name = "bright"
	p.SetMinValue(0)
	p.SetMaxValue(100)
	p.SetStepValue(1)
	p.SetValue(0)
	p.Unit = UnitPercentage

	return &BrightnessProperty{p}
}
