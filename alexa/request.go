package alexa

import (
	"encoding/json"
	"log"
)

type Request struct {
	Request json.RawMessage `json:"request"`
}

type RequestHeader struct {
	Type string `json:"type"`
}

// Type parses request headers from alexa json, returns request type like
// "IntentRequest", "LaunchRequest" or "" if it failed to unmarshal.
func (r *Request) Type() (ret string) {
	rh := new(RequestHeader)
	if err := json.Unmarshal(r.Request, rh); err != nil {
		log.Printf("failed to unmarshal RequestHeader: %v (body: %s)\n", err, r.Request)
		return
	}
	return rh.Type
}

type IntentRequest struct {
	RequestHeader
	Intent `json:"intent"`
}

type Intent struct {
	Name string `json:"name"`
}
