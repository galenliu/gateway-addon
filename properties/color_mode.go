package properties

const TypeColorModeProperty = "ColorModeProperty"

type ColorModeProperty struct {
	*StringProperty
}

func NewColorModeProperty() *ColorModeProperty {
	p := NewStringProperty(TypeColorModeProperty)
	p.Type = TypeString
	p.Name = ColorModel
	p.SetValue("color")
	p.Enum = []string{"color", "temperature"}

	return &ColorModeProperty{p}
}
