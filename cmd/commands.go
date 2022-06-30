package main

import (
	"github.com/tweekes0/pal-bot/internal/sounds"

	"github.com/bwmarrin/discordgo"
)

// Bot will join the voice channel that is specified in config file
func (app *application) joinVoice(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if app.joinedVoice {
		return ErrBotAlreadyJoinedVC
	}

	var err error
	app.vc, err = s.ChannelVoiceJoin(m.GuildID, app.botCfg.VoiceChannelID, false, true)
	if err != nil {
		return err
	}

	return nil
}

// Bot will leave the voice channel it is currently in
func (app *application) leaveVoice(s *discordgo.Session, m *discordgo.MessageCreate) error {
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
	if !app.joinedVoice {
		if err := app.joinVoice(s, m); err != nil {
			return err
		}
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
	for _, b := range buf {
		app.vc.OpusSend <- b
	}

	app.vc.Speaking(false)
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
	case startTime != "":
		start = startTime
		dur = 10
	case duration > 0 && duration <= 10:
		start = startTime
		dur = duration
	}

	f, err := sounds.CreateDCAFile(url, start, dur)
	if err != nil {
		return err
	}

	hash, err := sounds.HashFile(f)
	if err != nil {
		return err
	}

	_, err = app.soundbiteModel.Insert(name, m.Author.Username, m.Author.ID, f.Name(), hash)
	if err != nil {
		return err
	}

	return nil
}

// Bot will delete the specified sound 
func (app *application) deleteSound(s *discordgo.Session, m *discordgo.MessageCreate, name string) error {
	err := app.soundbiteModel.Delete(name, m.Author.ID)
	if err != nil {
		return err
	}

	return nil
}