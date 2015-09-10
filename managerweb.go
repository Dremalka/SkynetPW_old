package main

import (
	"github.com/labstack/echo"
	mid "github.com/labstack/echo/middleware"
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"
	"golang.org/x/net/websocket"
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

// websockDataBots метод отправляет по вебсокету к веб-клиенту обновленную информацию по ботам при получении сигнала из канала
func (mw *ManagerWeb) websockDataBots(c *echo.Context) error {
	act := Action{}		// структура действий с массивом активных вебсокетов (добавить, удалить...)
	ws := c.Socket()	// открытый вебсокет

	ch := make(chan Alarm)	// канал, по сигналу которого будет отправляться обновленная информация по боту на веб-клиенту через вебсокет

	act.Command = "add"	// добавить информацию по новому каналу и вебсокету
	act.Channel = ch
	mw.Sign <- act		// в массив активных вебсокетов

	defer func() {
		actdef := Action{}
		actdef.Command = "del"	// при закрытии вебсокета
		actdef.Channel = ch
		mw.Sign <- actdef			// удалить из массива вебсокет и канал
	}()

	type List struct {	// структура с данными, которые необходимо отправить на веб-клиент по вебсокету
		ID 		int 	`json:"id"`		// идентификатор бота
		Name	string 	`json:"name"`	// имя бота
	}

	ind := 0

	for {
		var st List	// создать структуру с данными бота
		st.ID	= ind
		st.Name	= "name-" + strconv.Itoa(ind)

		msg, _ := json.Marshal(st)	// сконвертировать структуру для отправки по вебсокету
		err := websocket.Message.Send(ws, string(msg))	// отправить данные веб-клиенту по вебсокету
		if err != nil {
			return err
		}

		<-ch	// ждать следующего сигнала для обновления информации
		ind++
	}
	return nil
}

// updateInfBots посылает сигнал всем обработчикам вебсокетов (websockDataBots) о том, что нужно обновить информацию на веб-клиентах
func (mw *ManagerWeb) updateInfBots(c *echo.Context) error {
	// перебрать массив каналов активных вебсокетов
	for i:=0; i<len(mw.Listch);i++ {
		ch := mw.Listch[i]
		ch <- Alarm{}	// каждому отправить сигнал, что необходимо обновить информацию на веб-клиенте
	}
	fmt.Println(mw.Listch)
	return c.String(http.StatusOK, "ok\n")
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
