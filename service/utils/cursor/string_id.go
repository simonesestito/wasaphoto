package cursor

import (
	"encoding/base64"
	"fmt"
	"github.com/gofrs/uuid"
	"strings"
)

func ParseStringIdCursor(userCursor string) (uuid.UUID, string, error) {
	if userCursor == "" {
		// Handle empty cursor
		return uuid.Nil, "", nil
	}

	rawCursorBytes, err := base64.URLEncoding.DecodeString(userCursor)
	if err != nil {
		return uuid.Nil, "", err
	}
	rawCursor := string(rawCursorBytes)

	// Split cursor
	cursorParts := strings.Split(rawCursor, ";")
	if len(cursorParts) != 2 {
		return uuid.Nil, "", fmt.Errorf("invalid cursor has %d parts (2 expected)", len(cursorParts))
	}
	rawId := cursorParts[0]
	rawSecondParam := cursorParts[1]

	// Parse ID
	id, err := uuid.FromString(rawId)
	if err != nil {
		return uuid.Nil, "", err
	}

	// Success!
	return id, rawSecondParam, nil
}

func CreateStringIdCursor(id []byte, secondParam string) string {
	rawId := uuid.FromBytesOrNil(id).String()
	rawCursorString := strings.Join([]string{rawId, secondParam}, ";")
	return base64.URLEncoding.EncodeToString([]byte(rawCursorString))
}
