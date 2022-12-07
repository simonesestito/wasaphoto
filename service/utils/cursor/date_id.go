package cursor

import (
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"time"
)

func ParseDateIdCursor(userCursor string) (uuid.UUID, time.Time, error) {
	// Parse as a string
	id, rawDate, err := ParseStringIdCursor(userCursor)
	if err != nil || userCursor == "" {
		return uuid.Nil, time.Now(), err
	}

	// Parse date
	date, err := timeprovider.UTCStringToDate(rawDate)
	if err != nil {
		return uuid.Nil, time.Now(), err
	}

	// Success!
	return id, date, nil
}

func CreateDateIdCursor(id []byte, date string) string {
	return CreateStringIdCursor(id, date)
}
