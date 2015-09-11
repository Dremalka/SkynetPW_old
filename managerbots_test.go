package main

import ("testing"
	"reflect"
)

func Test_NewManagerBots(t *testing.T) {
	mb, err := newManagerBots()
	if err != nil {
		t.Error("Типовая ошибка при создании нового объекта.")
	}
	if mb == nil {
		t.Error("метод newManagerBots() не должен возвращать nil.")
	} else {
		someval := make(map[string]*Bot)
		if reflect.DeepEqual(mb.ListBot, someval) != true {
			t.Errorf("Не инициализирован массив списка ботов. %v != %v", mb.ListBot, someval)
		}
	}


}

func Test_GetListBots(t *testing.T) {
	mb, err := newManagerBots()
	if err != nil {
		t.Error("Типовая ошибка при создании нового объекта.")
	}
	if mb == nil {
		t.Error("метод newManagerBots() не должен возвращать nil.")
	} else {
		// подготовить данные для тестирования
		bot := &Bot{}
		bot.ID = "12345"
		bot.Name = "name bot"
		bot.Server = "server"
		bot.Login = "login"
		bot.Password = "password"
		mb.ListBot[bot.ID] = bot

		// получить данные из проверяемого метода
		result, _ := mb.GetListBots()

		// проверить полученные данные
		if len(result) != 1 {
			t.Errorf("Неверное кол-во элементов в массиве. %d != %d\n", len(result), 1)
		} else {
			botRes := result[0]
			if botRes.ID != bot.ID {
				t.Errorf("Не совпадает ID заданного бота и проверяемого бота. %s != %s\n", bot.ID, botRes.ID)
			}
			if botRes.Name != bot.Name {
				t.Errorf("Не совпадает Name заданного бота и проверяемого бота. %s != %s\n", bot.Name, botRes.Name)
			}
			if botRes.Server != bot.Server {
				t.Errorf("Не совпадает Server заданного бота и проверяемого бота. %s != %s\n", bot.Server, botRes.Server)
			}
			if botRes.Login != bot.Login {
				t.Errorf("Не совпадает Login заданного бота и проверяемого бота. %s != %s\n", bot.Login, botRes.Login)
			}
			if botRes.Password != bot.Password {
				t.Errorf("Не совпадает Password заданного бота и проверяемого бота. %s != %s\n", bot.Password, botRes.Password)
			}

		}
	}

}
