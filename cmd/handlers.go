package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handler for when the bot receives a command
func (ctx *Context) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "" || !strings.HasPrefix(m.Content, ctx.botCfg.CommandPrefix) {
		return
	}

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
			return
		}
	} else {
		sl := strings.Split(c.command, ctx.botCfg.CommandPrefix)
		soundName := sl[len(sl)-1]
		exists, err := ctx.soundbiteModel.Exists(soundName, "")
		if err != nil {
			ctx.errorLogger.Println(err)
			return
		}

		if exists {
			err := ctx.playSound(s, m, soundName)
			if err != nil {
				ctx.errorLogger.Println(err)
				return
			}
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

	g, err := s.State.Guild(vs.GuildID)
	if err != nil {
		log.Fatal(err)
	}
	
	if len(g.VoiceStates) == 1 && ctx.joinedVoice {
		ctx.leaveVoice(s, nil)
	}
}
