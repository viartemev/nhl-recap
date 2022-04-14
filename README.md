# NHL recap
[![CI](https://github.com/viartemev/nhl-recap/actions/workflows/CI.yml/badge.svg?branch=master)](https://github.com/viartemev/nhl-recap/actions/workflows/CI.yml)

## Make commands
```shell
$make help

help                           This help.
docker_build_image             docker build
go_mod_verify                  go mod verify
go_build                       go build -v ./...
lint                           golint ./...
test                           go test -race -vet=off ./...

```

## Run bot
```shell
nhl_recap --help

Usage of nhl_recap:
  -t, --token string   Token for Telegram Bot API
```

## Run docker
```shell
docker run -d nhl_recap -t {TELEGRAM_BOT_TOKEN}
```
