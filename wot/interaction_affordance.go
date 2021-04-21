package wot

type InteractionAffordance struct {
	AtType       string       `json:"@type"`
	Title        string       `json:"title,omitempty"`
	Titles       string       `json:"titles,omitempty"`
	Description  string       `json:"description,omitempty"`
	Descriptions string       `json:"descriptions,omitempty"`
	Forms        []Form       `json:"forms,omitempty"`
	UriVariables []DataSchema `json:"uriVariables,omitempty"`
}
