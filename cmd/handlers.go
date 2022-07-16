package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handler for when the bot receives a command
func (ctx *Context) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	mention := fmt.Sprintf("<@%v>", m.Author.ID)

	if m.Author.ID == ctx.botID {
		return
	}

	if m.ChannelID != ctx.botCfg.BotChannelID && strings.HasPrefix(m.Content, ctx.botCfg.CommandPrefix) {
		msg := fmt.Sprintf("Psst.... %v I only respond to commands here", mention)
		_, _ = s.ChannelMessageSend(ctx.botCfg.BotChannelID, msg)
		return
	}
	
	if ctx.commands == nil {
		ctx.commands = ctx.getCommands(ctx.botCfg.CommandPrefix)
	} 

	c := parseCommand(m.Content)

	if command, ok := ctx.commands[c.command]; ok {
		err := command.Action(s, m, c.args)
		if err != nil {
			ctx.errorLogger.Println(err)
		}
	}
}

// Handler for when the bot joins or leaves a voice channel
func (ctx *Context) voiceStateChange(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.VoiceState.UserID == ctx.botID {
		ctx.isSpeaking = false

		if vs.VoiceState.ChannelID == "" { // The bot disconnects from a voice channel
			ctx.joinedVoice = false
		} else { // the bot joins a voice channel
			ctx.joinedVoice = true
		}
	}
}

func (ctx *Context) guildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	msg := "Hi @everyone, I am Pal Bot.\n**!commands** to see a list of my commands."
	_, _ = s.ChannelMessageSend(ctx.botCfg.BotChannelID, msg)
}
