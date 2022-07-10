package sounds

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/kkdai/youtube/v2"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

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

	return file, 0, nil
}

func createAACFile(url, startTime string, duration int) (*os.File, error) {
	videoFile, vidDuration, err := downloadYoutubeVideo(url)
	if err != nil {
		return nil, err
	}

	if startTime == "" && vidDuration > (10 * time.Second) {
		return nil, ErrTooLong
	}

	if duration > 10 {
		return nil, ErrTooLong
	}

	if startTime == "" {
		startTime = "00:00"
		duration = 10
	}

	name := getFilename(videoFile.Name())
	output := fmt.Sprintf("./audio/%v.aac", name)
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

	DeleteFile(videoFile.Name())

	return audio, nil
}

func CreateDCAFile(url, startTime string, duration int) (*os.File, *os.File, error) {
	aac, err := createAACFile(url, startTime, duration)
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
	f, err := os.Create(fmt.Sprintf("./audio/%v.dca", fname))
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

	return f, aac, nil
}
