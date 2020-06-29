package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/semyon-dev/znai-krai/config"
	defaultLog "log"
	"os"
	"strconv"
)

type ReportHook struct{}

var hooked = log.Hook(ReportHook{})

var bot *tgbotapi.BotAPI
var TelegramChatIDInt int

func Start() {
	// Создаём файл для логирования
	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		fmt.Println("✔ Логирование установлено в файл")
	} else {
		fmt.Println("× Failed to open file, using default stderr")
	}

	// Add file and line number to log
	log.Logger = log.With().Caller().Logger()

	// zerolog: создаём два канала для логирования - в консоль и файл
	multipleWriter := zerolog.MultiLevelWriter(os.Stdout, file)

	// дефолтное логирование направляем в zerolog
	defaultLog.SetFlags(0)
	defaultLog.SetOutput(multipleWriter)

	// Logging gin output to zerolog
	gin.DefaultWriter = multipleWriter
	gin.DefaultErrorWriter = multipleWriter

	connectBot()
}

func connectBot() {
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
	HandleErrWithMsg("", err)
}

func HandlePanicWitMsg(msg string, err error) {
	HandleErrWithMsg(msg, err)
	panic(err)
}

// логирование в консоль и телеграмм бота с сообщением
func HandleErrWithMsg(msg string, err error) {
	fmt.Println(msg, err)
	hooked.Error().Msg("Ошибка на бэкенде: " + msg + " " + err.Error())
}

func (h ReportHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.ErrorLevel {
		e.Str("info", level.String())
		sendToBot(msg)
	}
}

func sendToBot(text string) {
	msg := tgbotapi.NewMessage(int64(TelegramChatIDInt), text)
	msg.DisableWebPagePreview = true
	_, err := bot.Send(msg)
	if err != nil {
		log.Err(err)
	}
}
