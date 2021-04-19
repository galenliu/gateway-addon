package properties

const TypeColorModeProperty = "ColorModeProperty"

type ColorModeProperty struct {
	*StringProperty
}

func NewColorModeProperty() *ColorModeProperty {
	p := NewStringProperty(TypeColorModeProperty)
	p.SetType(TypeString)
	p.Name = ColorModel
	p.SetValue("color")
	p.Enum = []interface{}{"color", "temperature"}

	return &ColorModeProperty{p}
}
