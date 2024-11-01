package bot

import (
	"fmt"
	"tg-sunday-league/models"
)

type IMessageFormater interface {
	GameDetailsMessage(game *models.Game, players, absentees *[]models.User) string
	HelpMessage() string
	formatUserList(l *[]models.User) string
}

type MessageFormatter struct{}

func (m *MessageFormatter) GameDetailsMessage(game *models.Game, players, absentees *[]models.User) string {
	playerList := m.formatUserList(players)
	absenteesList := m.formatUserList(absentees)
	return fmt.Sprintf("Game on %s\nLocation: %s\nOpponent: %s\nPlayers: %s\nAbsentees: %s",
		game.Date.Format("2006-01-02 15:04"),
		game.Location,
		game.Opponent,
		playerList,
		absenteesList)
}

func (m *MessageFormatter) HelpMessage() string {

	helpText := `Bot Commands:

	`
	for _, c := range commands {
		helpText += c.Name + " - " + c.Description + "\n"
	}
	return helpText
}

func (m *MessageFormatter) formatUserList(l *[]models.User) string {
	userlist := ""
	copy_list := *l
	for i, player := range copy_list {
		userlist += fmt.Sprintf("%d. %s", i+1, player.Name)
		if player.HasPaid {
			userlist += " ✅"
		}
		userlist += "\n"
	}
	return userlist
}
