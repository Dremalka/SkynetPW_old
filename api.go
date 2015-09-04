package main

import (
	"github.com/labstack/echo"
	"net/http"
)

// api
// отправляет json-ответ с массивом текущих ботов
func listbot(c *echo.Context) error {
	//	TODO сформировать список ботов и отправить клиенту
	// временные данные
	type tmpstr struct {
		ID string
		Name string
	}
	list := make([]tmpstr, 2)
	list[0].ID = "123"
	list[0].Name = "abc"
	list[1].ID = "456"
	list[1].Name = "dfg"

	return c.JSON(http.StatusOK, list)
}
func createbot(c *echo.Context) error {
	//	TODO создать нового бота
	return c.String(http.StatusOK, "ok\n")
}
func sendactiontobot (c *echo.Context) error {
	//	TODO переслать команду боту (старт, стоп...)
	return c.String(http.StatusOK, "ok\n")
}
func deletebot(c *echo.Context) error {
	//	TODO удалить бота
	return c.String(http.StatusOK, "ok\n")
}


//websocket
func websockdatabots(c *echo.Context) error {
	//	TODO создать вебсокет для динамического обмена данными
	return nil
}


//отладочные функции
func updateinfbots(c *echo.Context) error {
	//	TODO иницировать обновление по вебсокету данных по ботам
	return c.String(http.StatusOK, "ok\n")
}
func testbot(c *echo.Context) error {
	//	TODO функция для отладки бота
	return c.String(http.StatusOK, "ok\n")
}