package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// Struct for all the config elements found in 'config.toml'
type BotConfig struct {
	DiscordToken       string `toml:"DiscordToken"`
	CommandPrefix      string `toml:"CommandPrefix"`
	BotChannelID       string `toml:BotChannelID`
}

// Reads the filename and Unmarshalls all of the entries into a *BotConfig struct
func ReadConfig(filename string) (*BotConfig, error) {
	file, err := ioutil.ReadFile(filename)
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

// Reads environment variables and returns a DSN that conforms to 
// the go mysql driver, https://github.com/go-sql-driver/mysql
func GetDSN() string {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	db   := os.Getenv("MYSQL_DB")

	return fmt.Sprintf("%v:%v@(%v)/%v?parseTime=true", user, pass, host, db)
}
