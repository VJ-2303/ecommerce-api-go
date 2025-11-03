package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// an custom type envelope is used to enclose the
// JSON output in a strutured way
type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if id < 1 || err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// writeJSON function decodes the given data into JSON
// and write into provided ResponseWriter with given status code
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	return nil
}
