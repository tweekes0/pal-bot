package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/tweekes0/pal-bot/config"
	"github.com/tweekes0/pal-bot/internal/sounds"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Description string
	Help        string
	Action      func(*discordgo.Session, *discordgo.MessageCreate, []string) error
}

type Commands map[string]Command

// Bot will join the voice channel that is specified in config file
func (ctx *Context) joinVoice(s *discordgo.Session, m *discordgo.MessageCreate) error {
	id := getChannelID(s, m)
	if id == "" {
		return ErrUserNotInVC
	}

	var err error
	ctx.vc, err = s.ChannelVoiceJoin(m.GuildID, id, false, true)
	if err != nil {
		return err
	}

	return nil
}

// Wrapper function for the 'join' command
func (ctx *Context) joinCommand() func(*discordgo.Session, *discordgo.MessageCreate, []string) error {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, st []string) error {
		if err := ctx.joinVoice(s, m); err != nil {
			return err
		}

		return nil
	}
}

// Bot will leave the voice channel it is currently in
func (ctx *Context) leaveVoice(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !ctx.joinedVoice {
		return ErrBotNotInVC
	}

	if err := ctx.vc.Disconnect(); err != nil {
		return err
	}

	return nil
}

// Wrapper function for the 'leave' command
func (ctx *Context) leaveCommand() func(*discordgo.Session, *discordgo.MessageCreate, []string) error {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, st []string) error {
		if err := ctx.leaveVoice(s, m); err != nil {
			return err
		}

		return nil
	}
}

// Bot will load an audio file from disc and play it in the voice channel specified in config file
func (ctx *Context) playSound(s *discordgo.Session, m *discordgo.MessageCreate, name string) error {
	if sound, ok := ctx.soundbiteCache[name]; ok {
		if err := ctx.streamSoundBite(s, m, sound); err != nil {
			return err
		}

		return nil
	}

	soundbite, err := ctx.soundbiteModel.Get(name)
	if err != nil {
		return err
	}
	ctx.soundbiteCache[name] = soundbite

	if err := ctx.streamSoundBite(s, m, soundbite); err != nil {
		return err
	}

	return nil
}

// Bot will create audio file from youtube video
func (ctx *Context) clip(s *discordgo.Session, m *discordgo.MessageCreate, name, url, startTime string, duration int) error {
	start, dur := getRuntime(startTime, duration)

	f, mp3, err := sounds.CreateDCAFile(config.AUDIO_DIR, url, start, dur)
	if err != nil {
		return err
	}

	hash, err := sounds.HashFile(f.Name())
	if err != nil {
		return err
	}

	_, err = ctx.soundbiteModel.Insert(name, m.Author.Username, m.Author.ID, f.Name(), hash)
	if err != nil {
		return err
	}

	ms := &discordgo.MessageSend{
		Content: fmt.Sprintf("Your clip is ready. Play it with **!%v**", name),
		Files:   []*discordgo.File{createDiscordFile(name, mp3)},
	}

	_, _ = s.ChannelMessageSendComplex(m.ChannelID, ms)

	err = sounds.DeleteFile(mp3.Name())
	if err != nil {
		return nil
	}

	return nil
}

// Wrapper function for the 'clip' command
func (ctx *Context) clipCommand() func(*discordgo.Session, *discordgo.MessageCreate, []string) error {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, st []string) error {
		args, err := parseClipCommand(st)
		if err != nil {
			return err
		}

		if err := ctx.clip(s, m, args.Name, args.Url, args.Start, args.Duration); err != nil {
			ctx.help(s, m, "clip")
			return err
		}

		return nil
	}
}

// Bot will delete the specified sound
func (ctx *Context) deleteSound(s *discordgo.Session, m *discordgo.MessageCreate, name string) error {
	sound, err := ctx.soundbiteModel.Get(name)
	if err != nil {
		return err
	}

	// remove item from cache if it is there.
	delete(ctx.soundbiteCache, name)

	err = ctx.soundbiteModel.Delete(name, m.Author.ID)
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

// Wrapper function for the 'delete' command
func (ctx *Context) deleteCommand() func(*discordgo.Session, *discordgo.MessageCreate, []string) error {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, st []string) error {
		if len(st) < 1 {
			ctx.help(s, m, "delete")
			return ErrNotEnoughArgs
		}

		name := st[0]
		if err := ctx.deleteSound(s, m, name); err != nil {
			return err
		}

		return nil
	}
}

// Bot will show all the sounds that available.
func (ctx *Context) showSounds(s *discordgo.Session, m *discordgo.MessageCreate) error {
	sounds, err := ctx.soundbiteModel.GetAll()
	if err != nil {
		return err
	}

	if len(sounds) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "there are no sounds :(")
		return nil
	}

	var b strings.Builder
	fmt.Fprint(&b, "**Available Sounds:** \n")
	for _, sound := range sounds {
		fmt.Fprintf(&b, "%v\n", sound.Name)
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, b.String())
	return nil
}

// Wrapper function for the 'sounds' command
func (ctx *Context) soundsCommand() func(*discordgo.Session, *discordgo.MessageCreate, []string) error {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, st []string) error {
		if err := ctx.showSounds(s, m); err != nil {
			return err
		}

		return nil
	}
}

// Function for the 'ping' command
func (ctx *Context) pingCommand() func(*discordgo.Session, *discordgo.MessageCreate, []string) error {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, st []string) error {
		mention := fmt.Sprintf("<@%v>", m.Author.ID)
		_, err := s.ChannelMessageSend(ctx.botCfg.BotChannelID, "Pong :D "+mention)

		if err != nil {
			return err
		}

		return nil
	}
}

// Bot will send the list commands and their descriptiions
func (ctx *Context) listCommands(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var sb strings.Builder
	keys := []string{}

	for k := range ctx.commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		sb.Write([]byte(fmt.Sprintf("**%v**:\t%v\n", k, ctx.commands[k].Description)))
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, sb.String())
	return nil
}

// Wrapper function for the 'commands' command
func (ctx *Context) commandsCommand() func(*discordgo.Session, *discordgo.MessageCreate, []string) error {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, st []string) error {
		return ctx.listCommands(s, m)
	}
}

// Bot will send the command's help info as a mention
func (ctx *Context) help(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	var sb strings.Builder
	sb.Write([]byte(ctx.commands[fmt.Sprint(ctx.botCfg.CommandPrefix, command)].Help))
	sb.Write([]byte(fmt.Sprintf("\n<@%v>", m.Author.ID)))
	_, _ = s.ChannelMessageSend(m.ChannelID, sb.String())
}

// Wrapper function for the 'help' command
func (ctx *Context) helpCommand() func(*discordgo.Session, *discordgo.MessageCreate, []string) error {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, st []string) error {
		if len(st) < 1 {
			return ctx.listCommands(s, m)
		}

		name := st[0]

		if command, ok := ctx.commands[fmt.Sprint(ctx.botCfg.CommandPrefix, name)]; ok {
			if command.Help != "" {
				ctx.help(s, m, name)
			}
		}

		return nil
	}
}

func (ctx *Context) upload(s *discordgo.Session, m *discordgo.MessageCreate, name string) error {
	if len(m.Attachments) == 0 {
		return ErrNoAttachments
	}

	url := m.Attachments[0].URL
	mp3, err := sounds.DownloadFileFromURL(name, url, 10)
	if err != nil {
		return err
	}
	defer sounds.DeleteFile(mp3.Name())
	mp3.Seek(0, io.SeekStart)

	f, err := sounds.MP3ToDCA(config.AUDIO_DIR, mp3)
	if err != nil {
		return err
	}

	hash, err := sounds.HashFile(f.Name())
	if err != nil {
		return err
	}

	_, err = ctx.soundbiteModel.Insert(name, m.Author.Username, m.Author.ID, f.Name(), hash)
	if err != nil {
		return err
	}

	ms := &discordgo.MessageSend{
		Content: fmt.Sprintf("Your clip is ready. Play it with **!%v**", name),
		Files:   []*discordgo.File{createDiscordFile(name, mp3)},
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, ms)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (ctx *Context) uploadCommand() func(*discordgo.Session, *discordgo.MessageCreate, []string) error {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, st []string) error {
		if len(st) < 1 {
			ctx.help(s, m, "upload")
			return ErrNotEnoughArgs
		}

		err := ctx.upload(s, m, st[0])
		if err != nil {
			return err
		}

		return nil
	}
}

func (ctx *Context) rename(s *discordgo.Session, m *discordgo.MessageCreate, oldName, newName string) error {
	err := ctx.soundbiteModel.UpdateName(oldName, newName)
	if err != nil {
		return err
	}

	sound, err := ctx.soundbiteModel.Get(newName)
	if err != nil {
		return err
	}

	delete(ctx.soundbiteCache, oldName)

	ctx.soundbiteCache[newName] = sound
	return nil
}

func (ctx *Context) renameCommand() func(*discordgo.Session, *discordgo.MessageCreate, []string) error {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, st []string) error {
		if len(st) < 2 {
			ctx.help(s, m, "rename")
			return ErrNotEnoughArgs
		}

		err := ctx.rename(s, m, st[0], st[1])
		if err != nil {
			return err
		}

		return nil
	}
}
