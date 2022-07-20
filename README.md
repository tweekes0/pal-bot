# Pal Bot 

[![Go Report Card](https://goreportcard.com/badge/github.com/tweekes0/pal-bot)](https://goreportcard.com/report/github.com/tweekes0/pal-bot)

Pal Bot is a discord bot that leverages
the [discordgo](https://github.com/bwmarrin/discordgo) library to create soundbites from Youtube videos and play them to your Discord server. This bot is heavily inspired by [Beginbot's](https://www.twitch.tv/beginbot) sound command that allows viewers to create sound bites and have them play on stream.

## Running Pal Bot

Before Pal Bot can be added to your Discord server, you must create your own configuration file. See [example_config.toml](https://github.com/tweekes0/pal-bot/blob/main/example_config.toml). Although Pal Bot is Dockerized, it is not a stateless application and will be error prone if deployed to serverless solution.

### Docker and Docker-Compose (Recommeded)

```
git clone https://github.com/tweekes0/pal-bot
cd pal-bot/
docker build -t pal-bot .
docker run -it -v "$(pwd)"/data/audio:/pal-bot/audio -v "$(pwd)"/data/db:/pal-bot/db pal-bot
```

or

```
git clone https://github.com/tweekes0/pal-bot
cd pal-bot/
docker-compose up
```

### Standard

Requires [ffmpeg](https://ffmpeg.org/) and [dca](https://github.com/bwmarrin/dca) to be installed.

```
sudo apt install ffmpeg // fairly large application
git clone https://github.com/tweekes0/pal-bot
go get -u github.com/bwmarrin/dca/cmd/dca@latest
go run ./cmd/
```

## Commands
Commands must be prefixed with prefix defined in 'config.toml'
| **Command**  | **Description**                                               |
| ------------ | ------------------------------------------------------------- |
| **clip**     | Take a youtube video and create a soundbite from it           |
| **commands** | List all available commands                                   |
| **delete**   | Delete a soundbite the user created                           |
| **help**     | Get help and usage for specified command                      |
| **join**     | Joins to the user's current VoiceChannel                      |
| **leave**    | Leaves the current VoiceChannel                               |
| **ping**     | Pong :D                                                       |
| **sounds**   | List all available sounds                                     |
| **!upload**  | Upload an mp3 and create a sounbite from it                   |
| **!rename**  | Rename soundbiets                                             |

## Examples

- #### Create a custom length soundbite and play it in your current VoiceChannel
```
!clip jigglypuff https://www.youtube.com/watch?v=d2NTtbusUso 00:06 8 
!jigglypuff
```

- #### Delete the jigglypuff soundbite
```
!delete jigglypuff
```

- #### Rename soundbite and play it
```
!rename jigglypuff jp
!jp
```

## Acknowledgements
Pal Bot has been brought to you in part by:

- [discordgo](https://github.com/bwmarrin/discordgo) - Go bindings for Discord
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - Sqlite3 driver for go using database/sql
- [filetype](https://github.com/h2non/filetype) - Fast, dependency-free Go package to infer binary file types based on the magic numbers header signature
- [go-mp3](https://github.com/hajimehoshi/go-mp3) - An MP3 decoder in pure Go
- [youtube](https://github.com/kkdai/youtube/) - Go package to download Youtube videos
- [ffmpeg-go](https://github.com/u2takey/ffmpeg-go) - Golang bindings for ffmpeg
- [dca](https://github.com/bwmarrin/dca) - Specification & Tool the Discord Audio (dca) file format