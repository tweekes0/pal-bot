package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tweekes0/pal-bot/internal/config"
	"github.com/tweekes0/pal-bot/internal/models"

	"github.com/bwmarrin/discordgo"
)

type application struct {
	botID          string
	botCfg         *config.BotConfig
	errorLogger    *log.Logger
	infoLogger     *log.Logger
	vc             *discordgo.VoiceConnection
	soundbiteModel *models.SoundbiteModel
	joinedVoice    bool
}

func main() {
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdin, "INFO: ", log.Ldate|log.Ltime)

	cfg, err := config.ReadConfig()
	if err != nil {
		errLog.Fatalln(err)
	}

	bot, err := initializeBot(cfg.DiscordToken)
	if err != nil {
		errLog.Fatalln(err)
	}

	botID, err := getBotID(bot)
	if err != nil {
		errLog.Fatalln(err)
	}

	db, err := openDB(cfg.DBConnectionString)
	if err != nil {
		errLog.Fatalln(err)
	}

	app := &application{
		joinedVoice: false,
		botID:       botID,
		botCfg:      cfg,
		errorLogger: errLog,
		infoLogger:  infoLog,
		soundbiteModel: &models.SoundbiteModel{DB: db},
	}

	bot.AddHandler(app.messageCreate)
	bot.AddHandler(app.voiceStateChange)

	infoLog.Println("Bot is now running. Press CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	bot.Close()
}
