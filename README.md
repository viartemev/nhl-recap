# NHL recap telegeram bot
[![CI](https://github.com/viartemev/nhl-recap/actions/workflows/CI.yml/badge.svg?branch=master)](https://github.com/viartemev/nhl-recap/actions/workflows/CI.yml)
[![CodeQL](https://github.com/viartemev/nhl-recap/actions/workflows/codeql-analysis.yml/badge.svg?branch=master)](https://github.com/viartemev/nhl-recap/actions/workflows/codeql-analysis.yml)

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
docker run -d viartemev/nhl-recap -t {TELEGRAM_BOT_TOKEN}
```

## Development
### Run tests in short mode:
```shell
go test -v -short
```
### Data race check:
```shell
go test -race
```
### Test coverage:
```shell
go test -coverprofile=cover.prof
go tool cover -html=cover.prof
```