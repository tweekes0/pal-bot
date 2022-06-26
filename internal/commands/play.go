package commands

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
)

func loadSound(filepath string) ([][]byte, error) {
	b := make([][]byte, 0)
	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	var opusLen int16

	for {
		err = binary.Read(file, binary.LittleEndian, &opusLen)
		if err == io.EOF || err == io.ErrUnexpectedEOF{
			err = file.Close()
			if err != nil {
				return nil, err
			}

			return b, nil
		}

		if err != nil {
			return nil, err
		}

		inBuf := make([]byte, opusLen)
		err = binary.Read(file, binary.LittleEndian, &inBuf)

		if err != nil {
			return nil, err
		}

		b = append(b, inBuf)
	}
}

func PlaySound(vc *discordgo.VoiceConnection, filepath string) error {
	buf, err := loadSound(filepath)	
	if err != nil {
		return err
	}

	vc.Speaking(true)
	for _, b := range buf {
		vc.OpusSend <- b
	}

	vc.Speaking(false)
	return nil
}
