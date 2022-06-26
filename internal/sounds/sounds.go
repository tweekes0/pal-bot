package sounds

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/kkdai/youtube/v2"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func downloadYoutubeVideo(url string) (*os.File, error) {
	client := &youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		return nil, err
	}

	format := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &format[0])
	if err != nil {
		return nil, err
	}

	file, err := os.CreateTemp("", "*.mp4")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func createAACFile(url, startTime, duration string) (*os.File, error) {
	videoFile, err := downloadYoutubeVideo(url)
	if err != nil {
		return nil, err
	}

	if startTime == "" {
		startTime = "00:00:00"
	}

	name := getFilename(videoFile.Name())
	output := fmt.Sprintf("./audio/%v.aac", name)
	kwargs := ffmpeg.KwArgs{"ss": startTime, "vn": "", "acodec": "copy"}

	if duration != "" {
		kwargs["t"] = duration
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

	deleteFile(videoFile)

	return audio, nil
}

func CreateDCAFile(url, startTime, duration string) error {
	aac, err := createAACFile(url, startTime, duration)
	if err != nil {
		return err
	}

	c1 := exec.Command("ffmpeg", "-i", aac.Name(), "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	c2 := exec.Command("dca")

	c2.Stdin, err = c1.StdoutPipe()
	if err != nil {
		return err
	}

	f, err := os.Create("./audio/test.dca")
	if err != nil {
		return err
	}

	c2.Stdout = f

	err = c2.Start()
	if err != nil {
		return err
	}

	err = c1.Run()
	if err != nil {
		return err
	}

	err = c2.Wait()
	if err != nil {
		return err
	}

	err = deleteFile(aac)
	if err != nil {
		return err
	}

	return nil
}
