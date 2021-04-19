package wot

type Form struct {
	Href                string      `json:"href,omitempty"`
	Rel                 string      `json:"rel"` //w3c定义中无这个，解决和wot兼容
	ContentType         string      `json:"contentType,omitempty"`
	ContentCoding       string      `json:"contentCoding,omitempty"`
	Security            interface{} `json:"security,omitempty"`
	Scopes              string      `json:"scopes,omitempty"`
	Response            interface{} `json:"response,omitempty"`
	AdditionalResponses interface{} `json:"additionalResponses,omitempty"`
	Subprotocol         string      `json:"subprotocol,omitempty"`
	Op                  interface{} `json:"op,omitempty"`
}
