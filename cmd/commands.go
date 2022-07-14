package main

import (
	"fmt"
	"strings"

	"github.com/tweekes0/pal-bot/internal/sounds"

	"github.com/bwmarrin/discordgo"
)

type Commands map[string]func([]string) error

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

// wrapper function for the 'join' command
func (app *application) joinCommand(s *discordgo.Session, m *discordgo.MessageCreate) func([]string) error {
	return func(st []string) error {
		if err := app.joinVoice(s, m); err != nil {
			return err
		}

		return nil
	}
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

// wrapper function for the 'leave' command
func (app *application) leaveCommand(s *discordgo.Session, m *discordgo.MessageCreate) func([]string) error {
	return func([]string) error {
		if err := app.leaveVoice(s, m); err != nil {
			return err
		}

		return nil
	}
}

// Bot will load an audio file from disc and play it in the voice channel specified in config file
func (app *application) playSound(s *discordgo.Session, m *discordgo.MessageCreate, name string) error {
	if sound, ok := app.soundbiteCache[name]; ok {
		if err := app.streamSoundBite(s, m, sound); err != nil {
			return err
		}

		return nil
	}

	soundbite, err := app.soundbiteModel.Get(name)
	if err != nil {
		return err
	}
	app.soundbiteCache[name] = soundbite


	if err := app.streamSoundBite(s, m, soundbite); err != nil {
		return err
	}

	return nil
}

// wrapper function for the 'play' command
func (app *application) playCommand(s *discordgo.Session, m *discordgo.MessageCreate) func([]string) error {
	return func(st []string) error {
		if len(st) < 1 {
			return ErrNotEnoughArgs
		}

		name := st[0]
		if err := app.playSound(s, m, name); err != nil {
			return err
		}

		return nil
	}
}

// Bot will create audio file from youtube video
func (app *application) clip(s *discordgo.Session, m *discordgo.MessageCreate, name, url, startTime string, duration int) error {
	start, dur := getRuntime(startTime, duration)

	f, mp3, err := sounds.CreateDCAFile(url, start, dur)
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
		Files:   []*discordgo.File{createDiscordFile(name, mp3)},
	}

	_, _ = s.ChannelMessageSendComplex(m.ChannelID, ms)

	err = sounds.DeleteFile(mp3.Name())
	if err != nil {
		return nil
	}

	return nil
}

// wrapper function for the 'clip' command
func (app *application) clipCommand(s *discordgo.Session, m *discordgo.MessageCreate) func([]string) error {
	return func(st []string) error {
		args, err := parseClipCommand(st)
		if err != nil {
			return err
		}

		if err := app.clip(s, m, args.Name, args.Url, args.Start, args.Duration); err != nil {
			return err
		}

		return nil
	}
}

// Bot will delete the specified sound
func (app *application) deleteSound(s *discordgo.Session, m *discordgo.MessageCreate, name string) error {
	sound, err := app.soundbiteModel.Get(name)
	if err != nil {
		return err
	}

	// remove item from cache if it is there.
	delete(app.soundbiteCache, name)

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

// wrapper function for the 'delete' command
func (app *application) deleteCommand(s *discordgo.Session, m *discordgo.MessageCreate) func([]string) error {
	return func(st []string) error {
		if len(st) < 1 {
			return ErrNotEnoughArgs
		}

		name := st[0]
		if err := app.deleteSound(s, m, name); err != nil {
			return err
		}

		return nil
	}
}

// Bot will show all the sounds that available.
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

// wrapper function for the 'sounds' command
func (app *application) soundsCommand(s *discordgo.Session, m *discordgo.MessageCreate) func([]string) error {
	return func(st []string) error {
		if err := app.showSounds(s, m); err != nil {
			return err
		}

		return nil
	}
}

func (app *application) pingCommand(s *discordgo.Session, m *discordgo.MessageCreate) func([]string) error {
	return func(st []string) error {
		mention := fmt.Sprintf("<@%v>", m.Author.ID)
		_, err := s.ChannelMessageSend(app.botCfg.BotChannelID, "Pong :D "+mention)

		if err != nil {
			return err
		}

		return nil
	}
}
