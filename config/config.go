package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

const (
	CONFIG_FILE         = "./config.toml"
	DB_DRIVER           = "sqlite3"
	DB_FILENAME         = "pal-bot.db"
	DB_DIR              = "./data/db"
	AUDIO_DIR           = "./data/audio"
	CLIP_MAX_DURATION   = 10 // Time in seconds for the maximum duration of a soundbite when using the 'clip' command
	UPLOAD_MAX_DURATION = 25 // Time in seconds for the maximum duration of a soundbite when using the 'upload' command
)

// Struct for all the config elements found in 'config.toml'
type BotConfig struct {
	DiscordToken  string `toml:"DiscordToken"`
	CommandPrefix string `toml:"CommandPrefix"`
	BotChannelID  string `toml:BotChannelID`
}

// Reads a config file and unmarshalls all of the entries into a *BotConfig struct
func ReadConfig() (*BotConfig, error) {
	file, err := ioutil.ReadFile(CONFIG_FILE)
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
