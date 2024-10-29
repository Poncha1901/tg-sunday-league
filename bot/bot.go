package bot

import (
	"strings"
	"tg-sunday-league/services"
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	TelegramBot *telebot.Bot
	GameService *services.GameService
}

func NewBot(token string, gameService *services.GameService) (*Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	b := &Bot{
		TelegramBot: bot,
		GameService: gameService,
	}

	b.setupHandlers()
	return b, nil
}

func (b *Bot) setupHandlers() {
	b.TelegramBot.Handle("/new", b.handleNewGame)
	b.TelegramBot.Handle("/register", b.handleRegisterPlayer)
	b.TelegramBot.Handle("/help", b.handleHelp)
}

// handleNewGame processes the command and arguments for creating a game
func (b *Bot) handleNewGame(m *telebot.Message) {
	argsStr := strings.TrimPrefix(m.Text, "/new ")
	argsStr = strings.Trim(argsStr, "()")
	args := strings.Split(argsStr, ",")
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	if len(args) < 4 {
		b.TelegramBot.Send(m.Chat, "Invalid format. Please use:\n/new (YYYY-MM-DD, HH:MM, Location, Opponent)")
		return
	}
	res, err := b.GameService.CreateNewGame(m.Chat.ID, args)

	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}
	b.TelegramBot.Send(m.Chat, res)
}

func (b *Bot) handleRegisterPlayer(m *telebot.Message) {
	argsStr := strings.TrimPrefix(m.Text, "/register ")
	playerName := strings.TrimSpace(argsStr)

	// Check if player name is provided
	if playerName == "" {
		b.TelegramBot.Send(m.Chat, "Invalid format. Please use:\n/register NameOfThePlayer")
		return
	}
	// Extract player ID from the message
	playerID := m.Sender.ID
	chatID := m.Chat.ID

	registration, err := b.GameService.RegisterPlayer(&chatID, &playerID, &playerName)

	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}
	b.TelegramBot.Send(m.Chat, registration)
}
