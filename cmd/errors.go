package main

import (
	"errors"
)

var (
	ErrDiscordSession     = errors.New("error creating Discord sesssion")
	ErrDiscordConnection  = errors.New("error opening connection")
	ErrBotAlreadyJoinedVC = errors.New("bot is already joined to voice")
	ErrBotNotInVC         = errors.New("bot is not joined to voice")
	ErrUserNotInVC        = errors.New("user is not joined to voice")
	ErrInvalidClipCommand = errors.New("clip needs at least a name and youtube link")
)
