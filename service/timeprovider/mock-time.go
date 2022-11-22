package timeprovider

import "time"

type MockTimeProvider struct {
	MockTime time.Time
}

func (provider MockTimeProvider) Now() time.Time {
	return provider.MockTime
}

func (provider MockTimeProvider) UTCString() string {
	return DateToUTCString(provider.Now())
}
