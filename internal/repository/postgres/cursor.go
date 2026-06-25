package postgres

import (
	"encoding/base64"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.999Z07:00"

func decodeCursor(encodedTime string) (time.Time, error) {
	b, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(timeFormat, string(b))
}

func encodeCursor(time time.Time) string {
	return base64.StdEncoding.EncodeToString([]byte(time.Format(timeFormat)))
}
