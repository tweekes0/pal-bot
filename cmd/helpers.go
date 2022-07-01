package main

import (
	"database/sql"
	"encoding/binary"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
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
			args:    nil,
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
		return nil, ErrDiscordConnection
	}

	bot.StateEnabled = true
	bot.State.TrackVoice = true
	bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates 

	err = bot.Open()
	if err != nil {
		return nil, ErrDiscordSession
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

type clipArgs struct {
	Name     string
	Url      string
	Start    string
	Duration int
}

func parseClipCommand(args []string) (*clipArgs, error) {
	var err error

	c := &clipArgs{}
	switch len(args) {
	case 0:
	case 1:
		return nil, ErrInvalidClipCommand
	case 2:
		c.Name = args[0]
		c.Url = args[1]
		c.Start = "00:00"
		c.Duration = 10
	case 3:
		c.Name = args[0]
		c.Url= args[1]
		c.Start = args[2]
		c.Duration = 10
	case 4:
		c.Name = args[0]
		c.Url = args[1]
		c.Start = args[2]
		c.Duration, err = strconv.Atoi(args[3])
		if err != nil {
			return nil, err
		}
	default:
	}

	return c, nil
}

func getChannelID(s *discordgo.Session, m *discordgo.MessageCreate) string {
	for _, guild := range s.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				return vs.ChannelID
			}
		}
	}

	return ""
}