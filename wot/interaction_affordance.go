package wot

type Form map[string]string

func NewForm(args ...string) Form {
	m := make(map[string]string, 0)
	for i, _ := range args {
		if i%2 == 0 {
			continue
		}
		m[args[i-1]] = args[i]
	}
	return m
}

type InteractionAffordance struct {
	AtType       string `json:"@type"`
	Title        string `json:"title,omitempty"`
	Titles       string `json:"titles,omitempty"`
	Description  string `json:"description,omitempty"`
	Descriptions string `json:"descriptions,omitempty"`

	Forms []Form `json:"forms,omitempty"`

	UriVariables []DataSchema `json:"uriVariables,omitempty"`
}
