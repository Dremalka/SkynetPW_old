package main

//Bot основная структура бота
type Bot struct {
	ID 			string
	Name 		string
	Server 		string
	Login		string
	Password	string
	status 		string
}

// newBot создать новый объект бота
func newBot(infbot map[string]string) (*Bot, error) {
	bot := &Bot{}
	bot.Name = infbot["name"]
	bot.Server = infbot["server"]
	bot.Login = infbot["login"]
	bot.Password = infbot["password"]

	return bot, nil
}
