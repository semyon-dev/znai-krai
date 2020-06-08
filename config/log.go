package config

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strconv"
)

type ReportHook struct{}

var bot *tgbotapi.BotAPI
var TelegramChatIDInt int

func ConnectBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(TelegramAPIToken)
	if err != nil {
		log.Err(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	TelegramChatIDInt, err = strconv.Atoi(TelegramChatID)
	if err != nil {
		log.Err(err)
	}
}

func (h ReportHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.ErrorLevel {
		e.Str("info", level.String())
		sendToBot(msg)
	}
}

func sendToBot(text string) {
	msg := tgbotapi.NewMessage(int64(TelegramChatIDInt), text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Err(err)
	}
}
