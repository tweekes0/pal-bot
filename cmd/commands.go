package main

import (
	"fmt"
	"strings"

	"github.com/tweekes0/pal-bot/internal/sounds"

	"github.com/bwmarrin/discordgo"
)

// Bot will join the voice channel that is specified in config file
func (app *application) joinVoice(s *discordgo.Session, m *discordgo.MessageCreate) error {
	id := getChannelID(s, m)
	if id == "" {
		return ErrUserNotInVC
	}

	var err error
	app.vc, err = s.ChannelVoiceJoin(m.GuildID, id, false, true)
	if err != nil {
		return err
	}

	return nil
}

// Bot will leave the voice channel it is currently in
func (app *application) leaveVoice() error {
	if !app.joinedVoice {
		return ErrBotNotInVC
	}

	if err := app.vc.Disconnect(); err != nil {
		return err
	}

	return nil
}

// Bot will load an audio file from disc and play it in the voice channel specified in config file
func (app *application) playSound(s *discordgo.Session, m *discordgo.MessageCreate, name string) error {
	if err := app.joinVoice(s, m); err != nil {
		return err
	}

	if app.isSpeaking {
		return nil
	}

	soundbite, err := app.soundbiteModel.Get(name)
	if err != nil {
		return err
	}

	buf, err := loadSound(soundbite.FilePath)
	if err != nil {
		return err
	}

	app.vc.Speaking(true)
	app.isSpeaking = true
	for _, b := range buf {
		app.vc.OpusSend <- b
	}

	app.vc.Speaking(false)
	app.isSpeaking = false
	return nil
}

// Bot will create audio file from youtube video
func (app *application) clip(s *discordgo.Session, m *discordgo.MessageCreate, name, url, startTime string, duration int) error {
	var start string
	var dur int

	switch {
	case startTime == "":
		start = ""
		dur = 10
	case startTime != "" && duration > 10:
		start = startTime
		dur = 10
	case duration > 0 && duration <= 10:
		start = startTime
		dur = duration
	}

	f, aac, err := sounds.CreateDCAFile(url, start, dur)
	if err != nil {
		return err
	}

	hash, err := sounds.HashFile(f.Name())
	if err != nil {
		return err
	}

	_, err = app.soundbiteModel.Insert(name, m.Author.Username, m.Author.ID, f.Name(), hash)
	if err != nil {
		return err
	}

	ms := &discordgo.MessageSend{
		Content: fmt.Sprintf("Your clip is ready. Play it with **!play %v**", name),
		Files: []*discordgo.File{createDiscordFile(name, aac)},
	}

	_,_ = s.ChannelMessageSendComplex(m.ChannelID, ms)

	err = sounds.DeleteFile(aac.Name())
	if err != nil {
		return nil
	}

	return nil
}

// Bot will delete the specified sound 
func (app *application) deleteSound(s *discordgo.Session, m *discordgo.MessageCreate, name string) error {
	sound, err := app.soundbiteModel.Get(name)
	if err != nil {
		return err
	}

	err = app.soundbiteModel.Delete(name, m.Author.ID)
	if err != nil {
		return err
	}

	err = sounds.DeleteFile(sound.FilePath)
	if err != nil { 
		return err
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v has been deleted\n", name))

	return nil
}

func (app *application) showSounds(s *discordgo.Session, m *discordgo.MessageCreate) error {
	sounds, err := app.soundbiteModel.GetAll()
	if err != nil {
		return err
	}

	if len(sounds) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "there are no sounds :(")
		return nil
	}

	var b strings.Builder
	fmt.Fprint(&b, "Current Sounds: \n")
	for _, sound := range sounds {
		fmt.Fprintf(&b, "%v\n", sound.Name)
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, b.String())
	return nil
}
