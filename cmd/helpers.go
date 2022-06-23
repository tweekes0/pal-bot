package main

import (
	"errors"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

type botConfig struct {
	DiscordToken  string `toml:"DiscordToken"`
	CommandPrefix string `toml:"CommandPrefix"`
}

func readConfig() (*botConfig, error) {
	file, err := ioutil.ReadFile("./config.toml")
	if err != nil {
		return nil, err
	}

	config := botConfig{}

	err = toml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

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