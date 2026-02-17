package main

import (
	"bytes"
	"encoding/json"
)

func prettyJSON(raw string) (string, error) {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(raw), "", "  ")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
