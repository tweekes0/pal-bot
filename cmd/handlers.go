package main

import (
	"fmt"
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
		args, err := parseClipCommand(c.args)
		if err != nil {
			app.errorLogger.Println(err)
			return
		}

		err = app.clip(s, m, args.Name, args.Url, args.Start, args.Duration)
		if err != nil {
			app.errorLogger.Println(err)
			return
		}

	case app.botCfg.CommandPrefix + "delete":
		if len(c.args) < 1 {
			return
		}

		name := c.args[0]
		if err := app.deleteSound(s, m, name); err != nil {
			app.errorLogger.Println(err)
			return
		}

	case app.botCfg.CommandPrefix + "sounds":
		if err := app.showSounds(s, m); err != nil {
			app.errorLogger.Println(err)
			return
		}

	default:
		return
	}
}

func (app *application) voiceStateChange(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.VoiceState.UserID == app.botID {
		app.isSpeaking = false

		if vs.VoiceState.ChannelID == "" { // The bot disconnects from a voice channel
			app.joinedVoice = false
		} else { // the bot joins a voice channel
			app.joinedVoice = true
		}
	}
}
