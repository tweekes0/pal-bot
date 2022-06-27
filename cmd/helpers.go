package main

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type botCommand struct {
	command string
	args    []string
}

func parseCommand(str string) *botCommand {
	s := strings.Fields(str)
	
	if len(s) == 0 {
		return nil
	}

	if len(s) == 1 {
		return &botCommand{
			command: s[0],
			args: nil,
		}
	}

	c := &botCommand{
		command: s[0],
	}
	for i := 1; i < len(s); i++ {
		c.args = append(c.args, s[i])
	}

	return c
}

func initializeBot(token string) (*discordgo.Session, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, errors.New(ErrDiscordConnection.Error())
	}

	bot.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

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

func loadSound(filepath string) ([][]byte, error) {
	b := make([][]byte, 0)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	var opusLen int16

	for {
		err = binary.Read(file, binary.LittleEndian, &opusLen)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = file.Close()
			if err != nil {
				return nil, err
			}

			return b, nil
		}

		if err != nil {
			return nil, err
		}

		inBuf := make([]byte, opusLen)
		err = binary.Read(file, binary.LittleEndian, &inBuf)
		if err != nil {
			return nil, err
		}

		b = append(b, inBuf)
	}
}
