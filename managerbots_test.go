package main

import (
	"testing"
	//"reflect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"fmt"
"errors"
)

func Test_NewManagerBots(t *testing.T) {
	mb, err := newManagerBots()
	assert.Nil(t, err, fmt.Sprintf("Типовая ошибка при создании нового объекта. Ошибка: %v", err))

	if result := assert.NotNil(t, mb, fmt.Sprint("Метод newManagerBots() не должен возвращать nil.")); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}

	someval := make(map[string]*Bot)
	assert.Equal(t, mb.ListBot, someval, "Не инициализирован массив списка ботов.")

	assert.NotNil(t, mb.db, "Не открыта база данных sqlite")
}

func Test_AddBot(t *testing.T) {

	mb, err := newManagerBots()
	if result := assert.Nil(t, err, fmt.Sprintf("Типовая ошибка при создании нового объекта. Ошибка: %v", err)); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}
	if result := assert.NotNil(t, mb, fmt.Sprint("Метод newManagerBots() не должен возвращать nil.")); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}


	infbot := make(map[string]string)
	infbot["name"] = "nametest"
	infbot["server"] = "127.0.0.1:29000"
	infbot["login"] = "logintest"
	infbot["password"] = "passwordtest"

	uid, err := mb.AddBot(infbot)
	assert.Nil(t, err, "Ошибка при создании бота.")
	if result := assert.NotEqual(t, uid, "", "Ошибка при создании бота. Возвращенный уид не должен быть пустым."); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}

	bot := mb.ListBot[uid]
	if result := assert.NotNil(t, bot, fmt.Sprintf("Ошибка при создании бота. В окружении mb.ListBot нет бота с уид-ом %s", uid)); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}

	assert.Equal(t, bot.Name, "nametest", "Имя бота не равно исходному.")
	assert.Equal(t, bot.Server, "127.0.0.1:29000", "Сервер не равен исходному.")
	assert.Equal(t, bot.Login, "logintest", "Логин не равен исходному.")
	assert.Equal(t, bot.Password, "passwordtest", "Пароль не равен исходному.")
}

func Test_SendActionToBot(t *testing.T) {
	// подготовка
	mb, err := newManagerBots()
	if result := assert.Nil(t, err, fmt.Sprintf("Типовая ошибка при создании нового объекта. Ошибка: %v", err)); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}
	if result := assert.NotNil(t, mb, fmt.Sprint("Метод newManagerBots() не должен возвращать nil.")); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}
	infbot := make(map[string]string)
	infbot["name"] = "nametest"
	infbot["server"] = "127.0.0.1:29000"
	infbot["login"] = "logintest"
	infbot["password"] = "passwordtest"

	uid, err := mb.AddBot(infbot)
	if result := assert.Nil(t, err, "Ошибка при создании бота."); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}
	if result := assert.NotEqual(t, uid, "", "Ошибка при создании бота. Возвращенный уид не должен быть пустым."); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}

	// тестирование
	result := mb.SendActionToBot("wrongid", "wrongaction", make(map[string]interface{}))
	temperr := errors.New("Не найден бот с идентификатором: wrongid")
	assert.Equal(t, result, temperr, "Метод должен вернуть ошибку.")



}

func Test_GetListBots(t *testing.T) {
	mb, err := newManagerBots()
	if result := assert.Nil(t, err, fmt.Sprintf("Типовая ошибка при создании нового объекта. Ошибка: %v", err)); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}
	if result := assert.NotNil(t, mb, fmt.Sprint("Метод newManagerBots() не должен возвращать nil.")); result != true {
		require.FailNow(t, "Дальнейшее тестирование функции прервано.")
	}


	// подготовить данные для тестирования
	bot := &Bot{}
	bot.ID = "12345"
	bot.Name = "name bot"
	bot.Server = "server"
	bot.Login = "login"
	bot.Password = "password"
	mb.ListBot[bot.ID] = bot

	// получить данные из проверяемого метода
	result, err := mb.GetListBots()

	assert.Nil(t, err, "Метод не должен возвращать ошибку.")
	assert.Equal(t, len(result), 1, "Неверное кол-во элементов в окружении ListBot.")

	botRes := result[0]
	assert.Equal(t, botRes.ID, bot.ID, "Не совпадает ID заданного бота и проверяемого бота.")
	assert.Equal(t, botRes.Name, bot.Name, "Не совпадает Name заданного бота и проверяемого бота.")
	assert.Equal(t, botRes.Server, bot.Server, "Не совпадает Server заданного бота и проверяемого бота.")
	assert.Equal(t, botRes.Login, bot.Login, "Не совпадает Login заданного бота и проверяемого бота.")
	assert.Equal(t, botRes.Password, bot.Password, "Не совпадает Password заданного бота и проверяемого бота.")
}
