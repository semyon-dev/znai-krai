package log

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/semyon-dev/znai-krai/config"
	"strconv"
)

type ReportHook struct{}

var hooked = log.Hook(ReportHook{})

var bot *tgbotapi.BotAPI
var TelegramChatIDInt int

func ConnectBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(config.TelegramAPIToken)
	if err != nil {
		log.Err(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	TelegramChatIDInt, err = strconv.Atoi(config.TelegramChatID)
	if err != nil {
		log.Err(err)
	}
}

// логирование в консоль и телеграмм бота
func HandleErr(err error) {
	fmt.Println(err)
	hooked.Error().Msg("Ошибка на бэкенде: " + err.Error())
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