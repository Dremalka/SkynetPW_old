package main

import (
	"github.com/labstack/echo"
	mid "github.com/labstack/echo/middleware"
	"net/http"
)

// ManagerWeb основная стркутура объекта Менеджер веб-интерфейса
type ManagerWeb struct {
	e *echo.Echo
//	TODO структура менеджера веб-интерфейса
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

		mw.e.Use(mid.Logger())	// выводить лог
		//mw.e.Use(mid.Recover())	// игнорировать ошибки при работе сервера

		mw.e.Get("/", hello)	// будущая основная страница

		//api
		mw.e.Get("/api/bots", listbot)	// вывести json-список текущих ботов
		mw.e.Post("/api/bots", createbot)	// создать нового бота
		mw.e.Patch("/api/bot/:id/:action", sendactiontobot)	// отправить основные команды боту (старт, стоп...)
		mw.e.Delete("/api/bot/:id", deletebot)	// удалить бота

		//websocket
		mw.e.WebSocket("/bots/ws", websockdatabots)	// вебсокет для динамического обновления информация по списку ботов

		// служебные вызовы на время разработки
		mw.e.Get("/api/bots/upd", updateinfbots)	// иницировать обновление информации в списке ботов
		mw.e.Post("/api/bot/test", testbot)

//		TODO инициализация настроек сервера (ip с которых можно принмать запросы, порт и т.д.)
		mw.e.Run(":8080")
	}()

	return nil
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
