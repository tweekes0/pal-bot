package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type application struct {
	botID       string
	botCfg   *botConfig
	errorLogger *log.Logger
	infoLogger  *log.Logger
}

func main() {
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdin, "INFO: ", log.Ldate|log.Ltime)

	cfg, err := readConfig()
	if err != nil {
		errLog.Println(err)
	}

	bot, err := initializeBot(cfg.DiscordToken)
	if err != nil {
		errLog.Println(err)
	}

	botID, err := getBotID(bot)
	if err != nil {
		errLog.Println(err)
	}

	app := &application{
		botID:       botID,
		botCfg:         cfg,
		errorLogger: errLog,
		infoLogger:  infoLog,
	}

	bot.AddHandler(app.defaultHandler)

	infoLog.Println("Bot is now running. Press CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	bot.Close()
}

func (app *application) defaultHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == app.botID {
		return
	}

	switch m.Content {
	case app.botCfg.CommandPrefix + "ping":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong :D")
	default:
		return
	}
}
