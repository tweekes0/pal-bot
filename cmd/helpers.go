package main

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)


// Struct to structure command received from *discordgo.Message.Content
type botCommand struct {
	command string
	args    []string
}

// Parses the specific command and any arguments that it may have 
func parseCommand(command string) *botCommand {
	s := strings.Fields(command)

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

// Creates a session to the discord API with a discord token.
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

// Gets the bot's ID, needed so the bot won't respond to itself.
func getBotID(bot *discordgo.Session) (string, error) {
	u, err := bot.User("@me")
	if err != nil {
		return "", err
	}

	return u.ID, nil
}

// Will load an DCA file into an 2d byte slice to then be played via opus connection
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

// Will make a connection to a mysql db with a given DSN
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

// Struct to structure arguments for 'clip' command
type clipArgs struct {
	Name     string
	Url      string
	Start    string
	Duration int
}

// Will parse the args from the 'clip' command and return a *clipArgs struct. 
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

// Gets the VoiceChannel of the user who sends a command, 
// will return nothing if the user is not in voice.
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

// Gets the startTime for a clip and its duration.
func getRuntime(start string, duration int) (startTime string, dur int) {
	switch {
	case start == "":
		startTime = ""
		dur = 10
	case start != "" && duration > 10:
		startTime = start
		dur = 10
	case duration > 0 && duration <= 10:
		startTime = start
		dur = duration
	}

	return
}

// Take a name and *os.File and transforms it into a discordgo.File
// needed for sending files to discord channel.
func createDiscordFile(name string, f *os.File) (*discordgo.File) {
	return &discordgo.File{
		Name: fmt.Sprintf("%v.mp3", name),
		Reader: f,
	}
}
