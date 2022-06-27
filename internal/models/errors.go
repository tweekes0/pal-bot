package models

import (
	"errors"
)

var (
	ErrNoRecords        = errors.New("no records were found")
	ErrCommandOwnership = errors.New("user did not create command")
)
