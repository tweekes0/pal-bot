package sounds

import (
	"errors"
)

var (
	ErrLengthTooLong    = errors.New("file is too long")
	ErrInvalidFile      = errors.New("file is not valid")
	ErrInvalidDuration  = errors.New("duration is not valid")
	ErrInvalidStartTime = errors.New("start time is not valid")
)
