package storage

import (
	"bytes"
	"encoding/json"
)

func jsonEscape(value string) string {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(value)
	if err != nil {
		return value
	}

	encoded := bytes.TrimSuffix(buf.Bytes(), []byte("\n"))

	return string(encoded[1 : len(encoded)-1])
}
