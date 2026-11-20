package storage

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"
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

func timestamp(t *time.Time) string {
	date := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0,
		0,
		0,
		0,
		t.Location(),
	)

	return strconv.FormatInt(date.Unix(), 10)
}

func field(t *time.Time) string {
	date := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		0,
		0,
		t.Location(),
	)

	return strconv.FormatInt(date.Unix(), 10)
}
