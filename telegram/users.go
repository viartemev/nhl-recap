package telegram

import "nhl-recap/util"

type TelegramUsers struct {
	Users *util.Set[int64]
}