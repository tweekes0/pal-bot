package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tweekes0/pal-bot/config"
	"github.com/tweekes0/pal-bot/internal/models"
	"github.com/tweekes0/pal-bot/internal/sounds"
	_ "modernc.org/sqlite"
)

// Descriptions for commands
const (
	pingDesc     = "Pong :D"
	joinDesc     = "Joins to the user's current VoiceChannel"
	leaveDesc    = "Leaves the current VoiceChannel"
	clipDesc     = "Take a youtube video and create a soundbite from it. Soundbites cannot be longer than 10 seconds.  **!help clip** for more info."
	playDesc     = "Play a sound that has been clipped. **!help play** for more info."
	deleteDesc   = "Delete a clipped soundbite the user created.  **!help delete** for more info."
	soundsDesc   = "List all available sounds"
	commandsDesc = "List all available commands"
	helpDesc     = "Get help and usage for specified commands"

	playHelp = `**!play** SOUNDNAME
**Example:** !play pika
Plays the 'pika' soundbite in the user's current VoiceChannel. Use **!sounds** to see all available sounds`
	clipHelp = `**!iclip** SOUNDNAME YOUTUBE_URL START_TIME(optional) DURATION(optional)
**Example:** !clip coolsound youtube.com/ID 00:23 5
Creates a new sound called 'coolsound' that starts at 00:23 and is 5 seconds long`
	deleteHelp = `**!delete** SOUNDNAME
**Example:** !delete pika
Deletes the soundbite the user created named 'pika'`
	helpHelp = `**!help** COMAND_NAME
**Example:** !help clip
Displays help information for the 'clip' command`
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

// Will make a connection to a sqlite db with a given filepath
func openDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open(config.DB_DRIVER, filepath)
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
		c.Url = args[1]
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

// Take a name and *os.File and transforms it into a discordgo.
// File needed for sending files to discord channel.
func createDiscordFile(name string, f *os.File) *discordgo.File {
	return &discordgo.File{
		Name:   fmt.Sprintf("%v.mp3", name),
		Reader: f,
	}
}

// Creates folder
func createFolder(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// Returns a map of all 'Commands'
func (app *application) getCommands(prefix string) Commands {
	c := make(Commands)
	c[fmt.Sprint(prefix, "ping")] = Command{
		description: pingDesc,
		help:        pingDesc,
		action:      app.pingCommand(),
	}
	c[fmt.Sprint(prefix, "join")] = Command{
		description: joinDesc,
		help:        joinDesc,
		action:      app.joinCommand(),
	}
	c[fmt.Sprint(prefix, "leave")] = Command{
		description: leaveDesc,
		help:        leaveDesc, 
		action:      app.leaveCommand(),
	}
	c[fmt.Sprint(prefix, "clip")] = Command{
		description: clipDesc,
		help:        clipHelp,
		action:      app.clipCommand(),
	}
	c[fmt.Sprint(prefix, "play")] = Command{
		description: playDesc,
		help:        playHelp,
		action:      app.playCommand(),
	}
	c[fmt.Sprint(prefix, "delete")] = Command{
		description: deleteDesc,
		help:        deleteHelp,
		action:      app.deleteCommand(),
	}
	c[fmt.Sprint(prefix, "sounds")] = Command{
		description: soundsDesc,
		help:        soundsDesc,
		action:      app.soundsCommand(),
	}
	c[fmt.Sprint(prefix, "commands")] = Command{
		description: commandsDesc,
		help:        commandsDesc,
		action:      app.commandsCommand(),
	}

	c[fmt.Sprint(prefix, "help")] = Command{
		description: helpDesc,
		help:        helpHelp,
		action:      app.helpCommand(),
	}

	return c
}

// Create a map for a cache of soundbites
func (app *application) createSoundsCache() (map[string]*models.Soundbite, error) {
	sounds, err := app.soundbiteModel.GetAll()
	if err != nil {
		return nil, err
	}

	cache := make(map[string]*models.Soundbite)
	for _, sound := range sounds {
		cache[sound.Name] = sound
	}

	return cache, nil
}

// Load and stream a soundbite into a VoiceChannel.
func (app *application) streamSoundBite(s *discordgo.Session, m *discordgo.MessageCreate, soundbite *models.Soundbite) error {
	if err := app.joinVoice(s, m); err != nil {
		return err
	}

	if app.isSpeaking {
		return nil
	}

	buf, err := sounds.LoadSound(soundbite.FilePath)
	if err != nil {
		return err
	}

	app.vc.Speaking(true)
	app.isSpeaking = true
	for _, b := range buf {
		app.vc.OpusSend <- b
	}

	app.vc.Speaking(false)
	app.isSpeaking = false
	return nil
}
