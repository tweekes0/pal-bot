package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type BotConfig struct {
	DiscordToken       string `toml:"DiscordToken"`
	CommandPrefix      string `toml:"CommandPrefix"`
	BotChannelID       string `toml:BotChannelID`
	DBConnectionString string `toml:DBConnectionString`
}

func ReadConfig() (*BotConfig, error) {
	file, err := ioutil.ReadFile("./config.toml")
	if err != nil {
		return nil, err
	}

	config := BotConfig{}

	err = toml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
