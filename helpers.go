package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
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
