package main

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mind1949/googletrans"
	"log"
	"os"
)

type RequestBody struct {
	Body string `json:"body"`
}

type Response struct {
	StatusCode int `json:"statusCode"`
}

func Handler(ctx context.Context, request RequestBody) (*Response, error) {

	var update tgbotapi.Update
	err := json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		log.Println(err)
		return &Response{StatusCode: 200}, nil
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	processUpdate(bot, &update)
	return &Response{StatusCode: 200}, nil

}

func processUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	if update.Message.IsCommand() {

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = "Привет, я переводчик!"
		case "help":
			msg.Text = "Отправь мне текст и я тебе его переведу, что может быть проще?"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}

	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		detected, err := googletrans.Detect(update.Message.Text)
		if err != nil {
			msg.Text = err.Error()
		} else {
			var dst string
			if detected.Lang == "ru" {
				dst = "en"
			} else {
				dst = "ru"
			}
			translate, err := Translate(update.Message.Text, dst)
			if err != nil {
				msg.Text = "Ошибка!\n" + translate
			} else {
				msg.Text = translate
			}
		}
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}

func Translate(str, dst string) (string, error) {

	params := googletrans.TranslateParams{
		Dest: dst,
		Text: str,
	}
	translated, err := googletrans.Translate(params)
	if err != nil {
		return err.Error(), err
	}
	return translated.Text, nil
}
