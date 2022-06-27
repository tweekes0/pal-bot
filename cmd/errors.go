package main

import (
	"errors"
)

var (
	ErrDiscordSession     = errors.New("error creating Discord sesssion")
	ErrDiscordConnection  = errors.New("error opening connection")
	ErrBotAlreadyJoinedVC = errors.New("error bot is already joined to voice")
	ErrBotNotInVC         = errors.New("error bot is not joined to voice")
)
