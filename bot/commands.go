package bot

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

type Command struct {
	Name        string
	Description string
}

var (
	HELP     = Command{"/help", "Show this help message"}
	NEW      = Command{"/new", "Create a new game with the specified date, time, location, and opponent.\n How to use: /new (YYYY-MM-DD, HH:MM, Location, Opponent)"}
	REGISTER = Command{"/register", "Register yourself with the specified name for the upcoming game"}
	DETAILS  = Command{"/details", "Show the details of the game"}
	PAID     = Command{"/paid", "Mark you as paid for the game"}
)
var commands = []Command{HELP, NEW, REGISTER, DETAILS, PAID}

type BotCommand interface {
	handleNewGame(m *telebot.Message)
	handleRegisterPlayer(m *telebot.Message)
	handleHelp(m *telebot.Message)
}

func (b *Bot) handleNewGame(m *telebot.Message) {
	argsStr := strings.TrimPrefix(m.Text, "/new ")
	argsStr = strings.Trim(argsStr, "()")
	args := strings.Split(argsStr, ",")
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	fmt.Println(len(args))
	if len(args) < 4 {
		b.TelegramBot.Send(m.Chat, "Invalid format. Please use:\n/new (YYYY-MM-DD, HH:MM, Location, Opponent, Optional[Price])")
		return
	}
	res, err := b.GameService.CreateNewGame(m.Chat.ID, m.Sender.ID, m.Sender.FirstName, args)

	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}
	b.TelegramBot.Send(m.Chat, res)
}

func (b *Bot) handleRegisterPlayer(m *telebot.Message) {
	playerID := m.Sender.ID
	playerName := m.Sender.FirstName
	log.Printf("Player ID: %d, Player Name: %s", playerID, playerName)
	chatID := m.Chat.ID

	registration, err := b.GameService.RegisterPlayer(&chatID, &playerID, &playerName)

	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}
	b.TelegramBot.Send(m.Chat, registration)
}

func (b *Bot) handleHelp(m *telebot.Message) {
	helpText := `Commands:`
	for _, c := range commands {
		helpText += "\n" + c.Name + " - " + c.Description
	}
	b.TelegramBot.Send(m.Chat, helpText)
}

func (b *Bot) handleDetails(m *telebot.Message) {
	game, err := b.GameService.GetGameDetails(m.Chat.ID)
	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}
	b.TelegramBot.Send(m.Chat, game)
}

func (b *Bot) handlePaid(m *telebot.Message) {
	playerID := m.Sender.ID
	playerName := m.Sender.FirstName
	chatID := m.Chat.ID

	paid, err := b.GameService.RepayGame(&chatID, &playerID, &playerName)

	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}
	b.TelegramBot.Send(m.Chat, paid)
}
