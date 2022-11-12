# syntax=docker/dockerfile:1
FROM golang:1.18 as build
WORKDIR /usr/local/go/src/pal-bot
COPY go.mod go.sum ./
RUN go mod download 
RUN go mod verify
RUN go install github.com/bwmarrin/dca/cmd/dca@latest
COPY . ./
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-w -s" -o /usr/local/bin/pal-bot ./cmd/

FROM ubuntu:latest 
RUN apt-get update \ 
    && apt-get install ca-certificates ffmpeg  youtube-dl -y \ 
    && update-ca-certificates 
WORKDIR /pal-bot/
COPY --from=build /usr/local/bin/pal-bot /usr/local/go/src/pal-bot/config.toml ./
COPY --from=build /go/bin/dca /usr/local/bin
ENTRYPOINT [ "./pal-bot" ]
