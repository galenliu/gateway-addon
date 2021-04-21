package properties

const TypeOnOffProperty = "OnOffProperty"

type OnOffProperty struct {
	*BooleanProperty
}

func NewOnOffProperty() *OnOffProperty {
	p := NewBooleanProperty(TypeOnOffProperty)
	p.SetType(TypeBoolean)
	p.Name = On
	p.SetValue(false)
	return &OnOffProperty{p}
}
