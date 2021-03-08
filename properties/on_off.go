package properties

const TypeOnOffProperty = "OnOffProperty"

type OnOffProperty struct {
	*BooleanProperty
}

func NewOnOffProperty() *OnOffProperty {
	p := NewBooleanProperty(TypeOnOffProperty)
	p.Type = TypeBoolean
	p.Name = On
	p.SetValue(false)
	return &OnOffProperty{p}
}
