package cursor

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"strings"
	"time"
)

func ParseDateIdCursor(userCursor string) (uuid.UUID, time.Time, error) {
	if userCursor == "" {
		// Handle empty cursor
		return uuid.Nil, time.Now(), nil
	}

	rawCursorBytes, err := base64.URLEncoding.DecodeString(userCursor)
	if err != nil {
		return uuid.Nil, time.Time{}, err
	}
	rawCursor := string(rawCursorBytes)

	// Split cursor
	cursorParts := strings.Split(rawCursor, ";")
	if len(cursorParts) != 2 {
		return uuid.Nil, time.Time{}, errors.New(fmt.Sprintf("invalid cursor has %d parts (2 expected)", len(cursorParts)))
	}
	rawId := cursorParts[0]
	rawDate := cursorParts[1]

	// Parse ID
	id, err := uuid.FromString(rawId)
	if err != nil {
		return uuid.Nil, time.Time{}, err
	}

	// Parse date
	date, err := timeprovider.UTCStringToDate(rawDate)
	if err != nil {
		return uuid.Nil, time.Time{}, err
	}

	// Success!
	return id, date, nil
}

func CreateDateIdCursor(id []byte, date string) string {
	rawId := uuid.FromBytesOrNil(id).String()
	rawCursorString := strings.Join([]string{rawId, date}, ";")
	return base64.URLEncoding.EncodeToString([]byte(rawCursorString))
}
