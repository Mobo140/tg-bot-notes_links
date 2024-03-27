package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnathorized  = errors.New("user is not authorize")
	errUnableToSave = errors.New("unable to save")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, "Произошла неизвестная ошибка.")

	switch err {
	case errUnathorized:
		msg.Text = "Ты не авторизирован! Используй команду /start ."
		b.bot.Send(msg)
	case errInvalidURL:
		msg.Text = "Это невалидная ссылка!"
		b.bot.Send(msg)
	case errUnableToSave:
		msg.Text = "Увы, не удалось сохранить ссылку. Попробуй повторить ещё раз."
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}

}
