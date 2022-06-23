package main

import (
	"errors"
)

var (
	ErrDiscordSession    = errors.New("error creating Discord sesssion")
	ErrDiscordConnection = errors.New("error opening connection")
)
