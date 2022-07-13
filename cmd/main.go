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
type application struct {
	botID          string
	botCfg         *config.BotConfig
	errorLogger    *log.Logger
	infoLogger     *log.Logger
	vc             *discordgo.VoiceConnection
	soundbiteModel *models.SoundbiteModel
	joinedVoice    bool
	isSpeaking     bool
}

func main() {
	// Loggers for outputting error and informational messages
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdin, "INFO: ", log.Ldate|log.Ltime)

	// Create audio folder 
	err := createFolder("./audio")
	if err != nil {
		errLog.Println(err)
	}

	// Create database folder
	err = createFolder("./db")
	if err != nil {
		errLog.Println(err)
	}

	// Create filepath and db connectivity
	path := filepath.Join(config.DB_DIR, config.DB_FILENAME)
	db, err := openDB(path)
	if err != nil {
		errLog.Fatalln(err)
	}

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

	app := &application{
		joinedVoice:    false,
		isSpeaking:     false,
		botID:          botID,
		botCfg:         cfg,
		errorLogger:    errLog,
		infoLogger:     infoLog,
		soundbiteModel: &models.SoundbiteModel{DB: db},
	}
	
	app.soundbiteModel.Initialize()

	bot.AddHandler(app.messageCreate)
	bot.AddHandler(app.voiceStateChange)

	infoLog.Println("Bot is now running. Press CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	bot.Close()
}
