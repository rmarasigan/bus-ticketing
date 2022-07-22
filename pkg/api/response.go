package api

import (
	"encoding/json"

	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
)

// EncodeResponse marshal a JSON encoding.
func EncodeResponse(body interface{}) []byte {
	encodeJSON, err := json.Marshal(body)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "EncodeResponseError", Message: "Failed to encode JSON"})
	}

	return encodeJSON
}

// ParseJSON unmarshal the JSON-encoded data and stores the result in the value pointed to by v.
func ParseJSON(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse JSON"})
		return err
	}

	return nil
}
