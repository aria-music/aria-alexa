package aria

import "encoding/json"

type ResponseHeader struct {
	Type string `json:"type"`
}

type RawResponse struct {
	ResponseHeader
	RawData json.RawMessage `json:"data"`
}

type Response struct {
	ResponseHeader
	Data interface{}
}

type StateData struct {
	Entry struct {
		Title string `json:"title"`
		URI   string `json:"uri"`
	} `json:"entry"`
}

var dateFactories = map[string]func() interface{}{
	"state": func() interface{} { return new(StateData) },
}
