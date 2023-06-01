package utility

import (
	"encoding/json"
)

// EncodeJSON marshals the body to produce a string format JSON.
func EncodeJSON(v any) string {
	encodeJSON, err := json.Marshal(v)
	if err != nil {
		Error(err, "JSONError", "failed to encode JSON")
	}

	return string(encodeJSON)
}

// ParseJSON unmarshal the JSON-encoded data and stores the result in the value pointed to by v.
func ParseJSON(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		Error(err, "JSONError", "failed to parse JSON")
		return err
	}

	return nil
}
