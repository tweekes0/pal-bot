package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/tweekes0/pal-bot/config"
	"github.com/tweekes0/pal-bot/internal/models"

	"github.com/bwmarrin/discordgo"
)

// Struct that holds the bot's loggers and state necessary
// to control the bot
type Context struct {
	botID          string
	botCfg         *config.BotConfig
	commands       Commands
	errorLogger    *log.Logger
	infoLogger     *log.Logger
	vc             *discordgo.VoiceConnection
	soundbiteModel *models.SoundbiteModel
	joinedVoice    bool
	isSpeaking     bool
	soundbiteCache map[string]*models.Soundbite
}

func main() {
	// Loggers for outputting error and informational messages
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdin, "INFO: ", log.Ldate|log.Ltime)

	if err := createFolder("./data"); err != nil {
		errLog.Fatalln(err)
	}

	if err := createFolder("./data/audio"); err != nil {
		errLog.Fatalln(err)
	}

	if err := createFolder("./data/db"); err != nil {
		errLog.Fatalln(err)
	}

	// Create filepath and db connectivity
	path := filepath.Join(config.DB_DIR, config.DB_FILENAME)
	db, err := openDB(path)
	if err != nil {
		errLog.Fatalln(err)
	}
	defer db.Close()

	// Create a struct from a config file
	cfg, err := config.ReadConfig()
	if err != nil {
		errLog.Fatalln(err)
	}

	// Create discord session from an API token
	bot, err := initializeBot(cfg.DiscordToken)
	if err != nil {
		errLog.Fatalln(err)
	}

	// Gets the discord bot id
	botID, err := getBotID(bot)
	if err != nil {
		errLog.Fatalln(err)
	}

	ctx := &Context{
		joinedVoice:    false,
		isSpeaking:     false,
		botID:          botID,
		botCfg:         cfg,
		errorLogger:    errLog,
		infoLogger:     infoLog,
		soundbiteModel: &models.SoundbiteModel{DB: db},
	}

	ctx.soundbiteModel.Initialize()

	// Create a cache of all the soundbites in the db
	soundbiteCache, err := ctx.createSoundsCache()
	if err != nil {
		errLog.Fatalln(err)
	}

	ctx.soundbiteCache = soundbiteCache

	bot.AddHandler(ctx.messageCreate)
	bot.AddHandler(ctx.voiceStateChange)
	bot.AddHandler(ctx.guildJoin)

	infoLog.Println("Bot is now running. Press CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	bot.Close()
}
