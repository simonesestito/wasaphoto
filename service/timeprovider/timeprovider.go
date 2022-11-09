package timeprovider

import "time"

type TimeProvider interface {
	Now() time.Time
}
