package sounds

import (
	"errors"
)

var (
	ErrInvalidDuration = errors.New("duration is not valid")
	ErrInvalidStartTime  =errors.New("start time is not valid")
)