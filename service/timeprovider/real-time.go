package timeprovider

import "time"

type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time {
	return time.Now()
}
