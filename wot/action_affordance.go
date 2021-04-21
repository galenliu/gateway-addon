package wot

type ActionAffordance struct {
	*InteractionAffordance

	Input      *DataSchema `json:"input,omitempty"`
	Output     *DataSchema `json:"output,omitempty"`
	Safe       bool        `json:"safe,omitempty"`
	Idempotent bool        `json:"idempotent,omitempty"`
}
