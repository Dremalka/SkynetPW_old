package main

import (
	"github.com/labstack/echo"
	mid "github.com/labstack/echo/middleware"
	"net/http"
	"fmt"
)

// ManagerWeb основная стркутура объекта Менеджер веб-интерфейса
type ManagerWeb struct {
	e *echo.Echo
	Sign chan Action	// канал по которому передается команда добавления/удаления вебсокета (нового подключения веб-клиента)
	Listch []chan Alarm // массив активных веб-клиентов
	//	TODO структура менеджера веб-интерфейса
}

// Action структура действий с массивом активных вебсокетов
type Action struct {
	Command string		// команда
	Channel chan Alarm	// канал, который слушает активный вебсокет. При получении сигнала обновляет информацию на веб-интерфейсе
}

// Alarm сигнал, который вызывает обновление информации в списке ботов и отправки ее по вебсокету
type Alarm struct {

}

func newManagerWeb() (*ManagerWeb, error) {
	mw := &ManagerWeb{}
	return mw, nil
}

//Start метод запускает веб-интерфейс
func (mw *ManagerWeb) Start() error {
	go func() {
		if mw.e == nil {
			mw.e = echo.New()
		}

		mw.e.Use(mid.Logger()) // выводить лог
		//mw.e.Use(mid.Recover())	// игнорировать ошибки при работе сервера

		mw.e.Get("/", hello) // будущая основная страница

		//api
		mw.e.Get("/api/bots", listbot)                      // вывести json-список текущих ботов
		mw.e.Post("/api/bots", createbot)                   // создать нового бота
		mw.e.Patch("/api/bot/:id/:action", sendactiontobot) // отправить основные команды боту (старт, стоп...)
		mw.e.Delete("/api/bot/:id", deletebot)              // удалить бота

		//websocket
		mw.e.WebSocket("/bots/ws", websockdatabots) // вебсокет для динамического обновления информация по списку ботов

		// служебные вызовы на время разработки
		mw.e.Get("/api/bots/upd", updateinfbots) // иницировать обновление информации в списке ботов
		mw.e.Post("/api/bot/test", testbot)

		//		TODO инициализация настроек сервера (ip с которых можно принимать запросы, порт и т.д.)
		mw.e.Run(":8080")
	}()

	return nil
}

// manager фоновая горутина, отслеживает актуальный список подключившихся веб-клиентов
func (mw *ManagerWeb) manager(ch <-chan Action) {
	for {
		act := <-ch // получить команду
		switch act.Command {
		case "add":		// добавить в масиив канал, который слушает новый вебсокет
			mw.Listch = append(mw.Listch, act.Channel)
		case "del":		// удалить из массива канал
			for i:= 0; i < len(mw.Listch); i++ {
				if mw.Listch[i] == act.Channel {
					mw.Listch = append(mw.Listch[:i], mw.Listch[i+1:]...)
					break
				}
			}
		}
		fmt.Println(mw.Listch)
	}
}

//Stop метод останавливает веб-интерфейс
func (mw *ManagerWeb) Stop() error {
	//	TODO остановка веб-интерфейса
	return nil
}

//Restart метод останавливает и запускает (перезапускает) веб-интерфейс
func (mw *ManagerWeb) Restart() error {
	var err error

	// остановить веб-интерфейс
	err = mw.Stop()
	if err != nil {
		return err
	}

	// запустить веб-интерфейс
	err = mw.Start()
	if err != nil {
		return err
	}
	return nil
}

func hello(c *echo.Context) error {
	//	TODO когда-нибудь будет выводить основную страницу
	return c.String(http.StatusOK, "ok\n")
}
