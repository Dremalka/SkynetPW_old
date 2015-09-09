package main
import (
	"fmt"
	"github.com/satori/go.uuid"
)

//ManagerBots основная структура менеджера ботов
type ManagerBots struct {
	ListBot map[string]*Bot
}

func newManagerBots() (*ManagerBots, error) {
	mb := &ManagerBots{}
	mb.ListBot = make(map[string]*Bot)
	return mb, nil
}

// AddBot метод создает нового бота и добавляет в массив
func (mb *ManagerBots) AddBot(infbot map[string]string) (string, error) {
	bot, err := newBot(infbot)
	if err != nil {
		fmt.Printf("Ошибка при создании нового бота. err: %v", err)
		return "", err
	}
	newuid := uuid.NewV4()
	mb.ListBot[newuid.String()] = bot
	return newuid.String(), nil
}

//Inf структура с информацией по боту для веб-интерфейса
type Inf struct {
	ID 			string
	Name 		string
	Server 		string
	Login 		string
	Password 	string
}

//GetListBots метод формирует массив с текущими ботами
func (mb *ManagerBots) GetListBots() ([]Inf, error) {
	var botinf Inf
	var list []Inf
	var id string
	var bot *Bot

	for id, bot = range mb.ListBot {
		botinf.ID = id
		botinf.Name = bot.Name
		botinf.Server = bot.Server
		botinf.Login = bot.Login
		botinf.Password = bot.Password
		list = append(list, botinf)
	}
	return list, nil
}
