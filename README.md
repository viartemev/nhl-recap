# NHL recap
```shell
$make help

help                           This help.
docker_build_image             docker build
go_mod_verify                  go mod verify
go_build                       go build -v ./...
lint                           golint ./...
test                           go test -race -vet=off ./...

```

```shell
nhl_recap --help

Usage of nhl_recap:
  -t, --token string   Token for Telegram Bot API
```

```shell
docker run -d nhl_recap -t {TELEGRAM_BOT_TOKEN}
```