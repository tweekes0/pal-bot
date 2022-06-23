package sounds

import (
	"fmt"
	"io"
	"os"

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

func createAudioFile(url, startTime, duration string) (*os.File, error) {
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

	fmt.Println(kwargs)

	err = ffmpeg.Input(videoFile.Name()).
		Output(output, kwargs).OverWriteOutput().Run()

	if err != nil {
		return nil, err
	}

	audio, err := os.Open(output)
	if err != nil {
		return nil, err
	}

	deleteVideoFile(videoFile)

	return audio, nil
}
