package sounds

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/kkdai/youtube/v2"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// Downloads a youtube video and returns an mp4 file if the download is succesful.
func downloadYoutubeVideo(url string) (*os.File, time.Duration, error) {
	client := &youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		return nil, 0, err
	}

	format := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &format[0])
	if err != nil {
		return nil, 0, err
	}

	file, err := os.CreateTemp("", "*.mp4")
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return nil, 0, err
	}

	return file, video.Duration, nil
}

// Converts an mp4 file to a AAC file.
func createAACFile(path, url, startTime string, duration int) (*os.File, error) {
	videoFile, d, err := downloadYoutubeVideo(url)
	if err != nil {
		return nil, err
	}

	st, err := startTimeToDuration(startTime)
	if err != nil {
		return nil, err
	}

	if st.Seconds() > d.Seconds() {
		return nil, ErrInvalidStartTime
	}

	if duration > 10 || duration < 1 {
		return nil, ErrInvalidDuration
	}

	if startTime == "" {
		startTime = "00:00"
		duration = 10
	}

	fname := getFilename(videoFile.Name())
	output := fmt.Sprintf("%v/%v.aac", path, fname)
	kwargs := ffmpeg.KwArgs{"ss": startTime, "vn": "", "acodec": "copy"}

	if duration != 0 {
		kwargs["t"] = fmt.Sprint(duration)
	}

	err = ffmpeg.Input(videoFile.Name()).
		Output(output, kwargs).OverWriteOutput().Run()

	if err != nil {
		return nil, err
	}

	audio, err := os.Open(output)
	if err != nil {
		return nil, err
	}
	defer audio.Close()

	DeleteFile(videoFile.Name())

	return audio, nil
}

// Converts and AAC file to a DCA file, file that can be streamed to discord VoiceChannel.
// Returns a the DCA file and an MP3 file that is needed to be sent as an embed to a TextChannel.
func CreateDCAFile(path, url, startTime string, duration int) (*os.File, *os.File, error) {
	aac, err := createAACFile(path, url, startTime, duration)
	if err != nil {
		return nil, nil, err
	}

	c1 := exec.Command("ffmpeg", "-i", aac.Name(), "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	c2 := exec.Command("dca")

	c2.Stdin, err = c1.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	fname := getFilename(aac.Name())
	f, err := os.Create(fmt.Sprintf("%v/%v.dca", path, fname))
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	c2.Stdout = f

	err = c2.Start()
	if err != nil {
		return nil, nil, err
	}

	err = c1.Run()
	if err != nil {
		return nil, nil, err
	}

	err = c2.Wait()
	if err != nil {
		return nil, nil, err
	}

	mp3, err := createMP3File(aac)
	if err != nil {
		return nil, nil, err
	}

	err = DeleteFile(aac.Name())
	if err != nil {
		return nil, nil, err
	}

	return f, mp3, nil
}

// Convert an AAC file into a MP3 using FFMPEG
func createMP3File(aac *os.File) (*os.File, error) {
	if aac == nil {
		return nil, ErrInvalidFile
	}

	mp3, err := os.CreateTemp("", "*.mp3")
	if err != nil {
		return nil, err
	}

	kwargs := ffmpeg.KwArgs{"acodec": "libmp3lame"}
	err = ffmpeg.Input(aac.Name()).
		Output(mp3.Name(), kwargs).OverWriteOutput().Run()

	if err != nil {
		return nil, err
	}

	return mp3, nil
}

// Will load an DCA file into an 2d byte slice to then be played via an opus connection
func LoadSound(filepath string) ([][]byte, error) {
	b := make([][]byte, 0)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	var opusLen int16

	for {
		err = binary.Read(file, binary.LittleEndian, &opusLen)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
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
