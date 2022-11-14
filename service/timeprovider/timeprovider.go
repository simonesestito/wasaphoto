package timeprovider

import "time"

const UTCFormat = "2006-01-02T15:04:05Z"

type TimeProvider interface {
	Now() time.Time
	UTCString() string
}
