package model

import (
	"encoding/json"
	"time"
)

type Duration time.Duration

// MarshalJSON implements json.Marshaller interface.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d))
}

// UnmarshalJSON implements json.Unmarshaller interface.
func (d *Duration) UnmarshalJSON(data []byte) (err error) {
	var duration string

	err = json.Unmarshal(data, &duration)
	if err != nil {
		return err
	}

	var parsed time.Duration

	parsed, err = time.ParseDuration(duration)
	if err != nil {
		return err
	}

	*d = Duration(parsed)

	return nil
}
