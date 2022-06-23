package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type application struct {
	discordBot  *discordgo.Session
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

	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord sesssion,", err)
		return
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		errLog.Println("error opening connection,", err)
		return
	}

	app := &application{
		discordBot:  dg,
		errorLogger: errLog,
		infoLogger:  infoLog,
	}

	fmt.Println(app)

	infoLog.Println("Bot is now running. Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()
}
