package data

import (
	"encoding/json"
	"time"
)

const myLayout = "2006-01-02 15:04:05"

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	tm := time.Time(t)

	if tm.IsZero() {
		return []byte("null"), nil
	}
	formattedString := tm.Format(myLayout)

	return json.Marshal(formattedString)
}
