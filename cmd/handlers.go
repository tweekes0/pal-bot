package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (app *application) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	mention := fmt.Sprintf("<@%v>", m.Author.ID)

	if m.Author.ID == app.botID {
		return
	}

	if m.ChannelID != app.botCfg.BotChannelID && strings.HasPrefix(m.Content, app.botCfg.CommandPrefix) {
		msg := fmt.Sprintf("Psst.... %v I only respond to commands here", mention)
		_, _ = s.ChannelMessageSend(app.botCfg.BotChannelID, msg)
		return
	}

	c := parseCommand(m.Content)

	switch c.command {
	case app.botCfg.CommandPrefix + "ping":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong :D "+mention)

	case app.botCfg.CommandPrefix + "join":
		if err := app.joinVoice(s, m); err != nil {
			app.errorLogger.Println(err)
		}

	case app.botCfg.CommandPrefix + "leave":
		if err := app.leaveVoice(s, m); err != nil {
			app.errorLogger.Println(err)
		}

	case app.botCfg.CommandPrefix + "play":
		if len(c.args) < 1 {
			return 
		}

		name := c.args[0]
		if err := app.playSound(s, m, name); err != nil {
			app.errorLogger.Println(err)
		}

	case app.botCfg.CommandPrefix + "clip":
		var (
			name     string
			url      string
			start    string
			duration int
			err      error
		)

		switch len(c.args) {
		case 0:
		case 1:
			msg := fmt.Sprintf(`clip needs a name and youtube link bud %v`, mention)
			_, _ = s.ChannelMessageSend(m.ChannelID, msg)
			return
		case 2:
			name = c.args[0]
			url = c.args[1]
			start = "00:00"
			duration = 10
		case 3:
			name = c.args[0]
			url = c.args[1]
			start = c.args[2]
			duration = 10
		case 4:
			name = c.args[0]
			url = c.args[1]
			start = c.args[2]
			duration, err = strconv.Atoi(c.args[2])
			if err != nil {
				app.errorLogger.Println(err)
				return
			}
		default:
		}

		err = app.clip(s, m, name, url, start, duration)
		if err != nil {
			app.errorLogger.Println(err)
			return
		}

	default:
		return
	}
}

func (app *application) voiceStateChange(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	app.joinedVoice = !app.joinedVoice
}
