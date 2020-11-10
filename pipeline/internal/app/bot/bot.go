package bot

import (
	"log"
	"strconv"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"gitlab.com/taskProvider/pipeline/internal/app/configure"
)

// Init ...
type Init struct {
	State bool
}

// NewBot ...
func NewBot(conf *configure.Conf, command, errors chan string, init *Init) {
	// register bot
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err)
	}

	// debug bot
	bot.Debug = conf.BotDebug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// update settings
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// error mesages
	go func() {
		for e := range errors {
			for _, item := range conf.BotAuth {
				num, _ := strconv.ParseInt(item, 10, 64)
				msg := tgbotapi.NewMessage(num, e)
				bot.Send(msg)
			}
		}
	}()

	// update mesages
	for update := range updates {
		// ignore any non-Message Updates
		if update.Message == nil { 
			continue
		}

		// get user id
		if update.Message.Text == "/userid" {
			st := strconv.Itoa(update.Message.From.ID)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, st))
			continue
		}

		// is user auth
		if botAuth(conf.BotAuth, update.Message.From.ID) == false {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извини бро, Я тебя незнаю...")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			continue
		}

		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(`
				Привет я бот для управления ` + conf.Name + `!
			Для получения доп. инфо отправь: /help .
			`))
			bot.Send(msg)
		}

		// get help
		if update.Message.Text == "/help" { 
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(`
				Привет я бот для управления микросервисами!

				Я умею распознавать команды:
				/off 	 - Выключить pipeline
				/on 	 - Включить pipeline
				/restart - Перезагрузить pipeline
				/logfile - Получить файл лога
				/errfile - Получить файл лога ошибок
				/userid  - получить useriD

				А так же я буду сообщать об ошибках сервера.
			`))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	
		// get log file
		if update.Message.Text == "/logfile" { 
			bot.Send(tgbotapi.NewDocumentUpload(update.Message.Chat.ID, conf.LogURL))
		}

		// get error log file
		if update.Message.Text == "/errfile" { 
			bot.Send(tgbotapi.NewDocumentUpload(update.Message.Chat.ID, conf.ErrorLogURL))
		}

		// off pipeline server
		if update.Message.Text == "/off" { 
			text := "Останавливаю  pipeline."
			if !init.State {
				text = "Pipeline уже выключен."
			} else {
				command <- "q"
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

		// on pipeline server
		if update.Message.Text == "/on" { 
			text := "Запускаю  pipeline."
			if init.State {
				text = "Pipeline уже запущен"
			} else {
				command <- "s"
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

		// restart pipeline server
		if update.Message.Text == "/restart" { 
			text := "Перезапускаю  pipeline."
			if !init.State {
				text = "Pipeline выключен, restart не выполнен"
			} else {
				command <- "r"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

// auth
func botAuth(auth []string, id int) bool {
	result := false;
	for _, item := range auth {
		if item == strconv.Itoa(id) {
			result = true
		}
	}

	return result
}