package main

import "github.com/bwmarrin/discordgo"

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

	app.joinedVoice = true
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

	app.joinedVoice = false

	return nil
}

// Bot will load an audio file from disc and play it in the voice channel specified in config file
func (app *application) playSound(s *discordgo.Session, m *discordgo.MessageCreate, filepath string) error {
	if !app.joinedVoice {
		if err := app.joinVoice(s, m); err != nil {
			return err
		}
	}

	buf, err := loadSound(filepath)
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
