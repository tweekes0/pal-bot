package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handler for when the bot receives a command
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
	
	if app.commands == nil {
		app.commands = app.getCommands(app.botCfg.CommandPrefix)
	} 

	c := parseCommand(m.Content)

	if command, ok := app.commands[c.command]; ok {
		err := command.action(s, m, c.args)
		if err != nil {
			app.errorLogger.Println(err)
		}
	}
}

// Handler for when the bot joins or leaves a voice channel
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

func (app *application) guildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	msg := "Hi @everyone, I am Pal Bot.\n**!commands** to see a list of my commands."
	_, _ = s.ChannelMessageSend(app.botCfg.BotChannelID, msg)
}
