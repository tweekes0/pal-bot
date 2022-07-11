package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type BotConfig struct {
	DiscordToken       string `toml:"DiscordToken"`
	CommandPrefix      string `toml:"CommandPrefix"`
	BotChannelID       string `toml:BotChannelID`
	// DBConnectionString string `toml:DBConnectionString`
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

func GetDSN() string {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	db := os.Getenv("MYSQL_DB")

	return fmt.Sprintf("%v:%v@(%v)/%v?parseTime=true", user, pass, host, db)

}