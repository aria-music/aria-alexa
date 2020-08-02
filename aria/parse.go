package aria

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func ParseResponse(body io.Reader) (*Response, error) {
	rr, err := ParseRawResponse(body)
	if err != nil {
		// TODO: will be fixed in aria-core
		log.Printf("failed to get RawResponse. Returning empty response...")
		return new(Response), nil
		// return nil, fmt.Errorf("failed to get RawResponse: %w", err)
	}

	r, err := rr.ParseRawResponse()
	if err != nil {
		return nil, fmt.Errorf("failed to get Response from RawResponse: %w", err)
	}

	return r, nil
}

func ParseRawResponse(body io.Reader) (*RawResponse, error) {
	dec := json.NewDecoder(body)
	rr := new(RawResponse)
	if err := dec.Decode(rr); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return rr, nil
}

func (rr *RawResponse) ParseRawResponse() (*Response, error) {
	r := new(Response)
	r.ResponseHeader = rr.ResponseHeader

	dh, ok := dateFactories[r.Type]
	if !ok {
		return nil, fmt.Errorf("no data factory is found for response type (%s)", r.Type)
	}
	r.Data = dh()

	if err := json.Unmarshal(rr.RawData, r.Data); err != nil {
		return nil, fmt.Errorf("failed to decode RawData to json: %w", err)
	}

	return r, nil
}
