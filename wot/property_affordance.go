package wot

type PropertyAffordance struct {
	*InteractionAffordance
	IDataSchema
	Observable bool `json:"observable"`
}
