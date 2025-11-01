package main

import (
	"encoding/json"
	"net/http"
)

// an custom type envelope is used to enclose the
// JSON output in a strutured way
type envelope map[string]any

// writeJSON function decodes the given data into JSON
// and write into provided ResponseWriter with given status code
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	return nil
}
