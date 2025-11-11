package validator

import "regexp"

// Validator used to hold the validation Errors in a key - message format
// Useful methods are implemented on it
type Validator struct {
	Errors map[string]string
}

var PhoneNumberRegex = regexp.MustCompile("^[0-9]{10}$")

// New creates and empty Validator and return a ready to use validator type
func New() *Validator {
	return &Validator{make(map[string]string)}
}

// Valid checks if the validator doesn't hold any Errors
// False means there are validation errors available and cant proceed further
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds the provided key, message in the validator's
// Errors fields, if only its not exists
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check call the AddError using the provided key, message
// if given condition in ok is not true
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
