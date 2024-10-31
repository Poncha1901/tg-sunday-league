package bot

import (
	"tg-sunday-league/services"
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	TelegramBot     *telebot.Bot
	MessageFormater IMessageFormater
	GameService     services.IGameService
}

func NewBot(token string, gameService services.IGameService, messageFormater IMessageFormater) (*Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	b := &Bot{
		TelegramBot:     bot,
		MessageFormater: messageFormater,
		GameService:     gameService,
	}

	b.setupHandlers()
	return b, nil
}

func (b *Bot) setupHandlers() {
	b.TelegramBot.Handle(NEW.Name, b.handleNewGame)
	b.TelegramBot.Handle(REGISTER.Name, b.handleRegisterPlayer)
	b.TelegramBot.Handle(OUT.Name, b.handleRegisterPlayer)
	b.TelegramBot.Handle(HELP.Name, b.handleHelp)
	b.TelegramBot.Handle(DETAILS.Name, b.handleDetails)
	b.TelegramBot.Handle(PAID.Name, b.handlePaid)
}
