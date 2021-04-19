package wot

type EventAffordance struct {
	*InteractionAffordance
	Subscription *DataSchema `json:"subscription,omitempty"`
	Data         *DataSchema `json:"data,omitempty"`
	Cancellation *DataSchema `json:"cancellation,omitempty"`
}
