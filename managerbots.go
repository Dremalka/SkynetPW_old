package main
import (
	"fmt"
	"github.com/satori/go.uuid"
	"errors"
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

// SendActionToBot выполнить действие
func (mb *ManagerBots) SendActionToBot(id, action string, param map[string]interface{}) error {
	var err error
	bot := mb.ListBot[id]
	if bot == nil {
		fmt.Printf("Не найден бот с идентификатором: %s\n", id)
		err = errors.New("Не найден бот с идентификатором: " + id)
		return err
	}

	fmt.Println("action", action)

	switch action {
	case "start":	// запустить бота
		err = bot.Start()
		if err != nil {
			return err
		}
	case "stop":	// остановить бота
		err = bot.Stop()
		if err != nil {
			return err
		}
	case "exit":	//удалить бота
		err = bot.Exit()
		if err != nil {
			return err
		}
		delete(mb.ListBot, id)
	case "connect":	// подключить бота к серверу
//		TODO подключиться к серверу
	case "disconnect":	// отключить бота от сервера
//	TODO отключиться от сервера
	}
	return nil
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
