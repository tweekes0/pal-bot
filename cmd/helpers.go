package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func initializeBot(token string) (*discordgo.Session, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, errors.New(ErrDiscordConnection.Error())
	}

	bot.Identify.Intents = discordgo.IntentsGuildMessages
	err = bot.Open()
	if err != nil {
		return nil, errors.New(ErrDiscordSession.Error())
	}

	return bot, nil
}

func getBotID(bot *discordgo.Session) (string, error) {
	u, err := bot.User("@me")
	if err != nil {
		return "", err
	}

	return u.ID, nil
}
