package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// an custom type envelope is used to enclose the
// JSON output in a strutured way
type envelope map[string]any

// readIDParam params get all the url parameters from the request context
// And get the specific paramter named 'id', And convert it to integer before returning
// if the conversion fails or the 'id' is less than 1 it returns an error message
func (app *application) readIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if id < 1 || err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// writeJSON function decodes the given data into JSON
// and write into provided ResponseWriter with given status code
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope) error {
	// Marshal the given struct into and JSON bytes array
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	// Add an newline character at the end
	// for easier readind and commandLine outputs
	js = append(js, '\n')

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON bytes array to the response body
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Only Maximum of 1mb of request body is allowed
	// anything more than that will return an error
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	// Strictly disallow unknown fields which are not in the struct
	// pointed by the dst, if unknown fields available it
	// will produce an error
	dec.DisallowUnknownFields()

	// decode the JSON body into the struct pointed by dst
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// A json.InvalidUnmarshalError error will be returned if we pass something
		// that is not a non-nil pointer to Decode(). We catch this and panic,
		// rather than returning an error to our handler. At the end of this chapter
		// we'll talk about panicking versus returning errors, and discuss why it's an
		// appropriate thing to do in this specific situation.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be loarger than %d bytes", maxBytesError.Limit)

		default:
			return err
		}
	}

	// Decode the request body second time, and check if its returning the io.EOF error
	// if Not, it means the body have more than one JSON body
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
