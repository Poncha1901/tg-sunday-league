package bot

import (
	"fmt"
	"log"
	"strings"
	"tg-sunday-league/services"

	"gopkg.in/tucnak/telebot.v2"
)

type Command struct {
	Name        string
	Description string
}

var (
	HELP = Command{"/help", `Show this help message`}
	NEW  = Command{"/new", `Create a new game with the specified date, time, location, opponent and price.
							How to use: /new (YYYY-MM-DD, HH:MM, Location, Opponent, Price)
							i.e: /new (2024-10-10, 11:00, Marina Bay Sands, Célavi FC, 15)`}
	CANCEL  = Command{"/cancel", `Cancel the upcoming game`}
	IN      = Command{"/in", `Register yourself for the upcoming game`}
	OUT     = Command{"/out", `Mark yourself as absent for the upcoming game`}
	DETAILS = Command{"/details", `Show the details of the game`}
	PAID    = Command{"/paid", `Mark you as paid for the game`}
)
var commands = []Command{HELP, NEW, IN, OUT, DETAILS, PAID}

type IBotCommand interface {
	handleNewGame(m *telebot.Message)
	handleRegisterPlayer(m *telebot.Message)
	handleHelp(m *telebot.Message)
	handleDetails(m *telebot.Message)
	handlePaid(m *telebot.Message)
	handleCancelGame(m *telebot.Message)
	isAdmin(bot *telebot.Bot, chat *telebot.Chat, user *telebot.User) bool
	isMessageSentFromGroup(m *telebot.Message) bool
}

func (b *Bot) handleNewGame(m *telebot.Message) {
	if !b.isMessageSentFromGroup(m) {
		return
	}
	if !b.isAdmin(m.Chat, m.Sender) {
		return
	}

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
	game, players, absentees, err := b.GameService.CreateNewGame(m.Chat.ID, m.Sender.ID, m.Sender.FirstName, args)
	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}
	message := b.MessageFormater.GameDetailsMessage(game, players, absentees)
	b.TelegramBot.Send(m.Chat, message)
}

func (b *Bot) handleCancelGame(m *telebot.Message) {
	if !b.isMessageSentFromGroup(m) {
		return
	}
	if !b.isAdmin(m.Chat, m.Sender) {
		return
	}

	game, err := b.GameService.CancelGame(m.Chat.ID)
	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}
	message := fmt.Sprintf("Game on %s has been cancelled.", game.Date.Format("2006-01-02 15:04"))
	b.TelegramBot.Send(m.Chat, message)
}

func (b *Bot) handleRegisterPlayer(m *telebot.Message) {
	if !b.isMessageSentFromGroup(m) {
		return
	}
	var status services.PlayerStatus
	log.Printf("Message: %s", m.Text)
	switch m.Text {
	case "/register":
		status = services.ATTENDING
	case "/out":
		status = services.OUT
	}
	playerID := m.Sender.ID
	var playerName string
	if m.Sender.FirstName == "" {
		playerName = m.Sender.Username
	} else {
		playerName = m.Sender.FirstName
	}

	log.Printf("Player ID: %d, Player Name: %s", playerID, playerName)
	chatID := m.Chat.ID

	game, players, absentees, err := b.GameService.RegisterPlayer(&chatID, &playerID, &playerName, status)
	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}

	message := b.MessageFormater.GameDetailsMessage(game, players, absentees)
	b.TelegramBot.Send(m.Chat, message)
}

func (b *Bot) handleHelp(m *telebot.Message) {
	b.TelegramBot.Send(m.Chat, b.MessageFormater.HelpMessage())
}

func (b *Bot) handleDetails(m *telebot.Message) {
	game, players, absentees, err := b.GameService.GetGameDetails(m.Chat.ID)
	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}
	message := b.MessageFormater.GameDetailsMessage(game, players, absentees)
	b.TelegramBot.Send(m.Chat, message)
}

func (b *Bot) handlePaid(m *telebot.Message) {
	if !b.isMessageSentFromGroup(m) {
		return
	}
	playerID := m.Sender.ID
	chatID := m.Chat.ID

	game, players, absentees, err := b.GameService.RepayGame(&chatID, &playerID)

	if err != nil {
		b.TelegramBot.Send(m.Chat, err.Error())
		return
	}

	message := b.MessageFormater.GameDetailsMessage(game, players, absentees)
	b.TelegramBot.Send(m.Chat, message)
}

func (b *Bot) isAdmin(chat *telebot.Chat, user *telebot.User) bool {
	admins, err := b.TelegramBot.AdminsOf(chat)
	if err != nil {
		log.Printf("Could not get admins of chat: %v", err)
		return false
	}
	for _, admin := range admins {
		if admin.User.ID == user.ID {
			return true
		}
	}
	b.TelegramBot.Send(chat, "Only admins of the group can create a new game.")
	return false
}

func (b *Bot) isMessageSentFromGroup(m *telebot.Message) bool {
	if !m.FromGroup() {
		b.TelegramBot.Send(m.Chat, "This bot is intended to work for group chat only. /help for more info")
		return false
	}
	return true
}
