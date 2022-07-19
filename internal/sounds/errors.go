package sounds

import (
	"errors"
)

var (
	ErrInvalidFile = errors.New("file is not valid")
	ErrInvalidDuration = errors.New("duration is not valid")
	ErrInvalidStartTime  =errors.New("start time is not valid")
)