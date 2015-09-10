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
	list, err := MB.GetListBots()
	if err != nil {
		return c.JSON(http.StatusConflict, nil)
	}
	return c.JSON(http.StatusOK, list)
}

func createbot(c *echo.Context) error {
	//	TODO создать нового бота
	return c.String(http.StatusOK, "ok\n")
}
func sendactiontobot(c *echo.Context) error {
	//	TODO переслать команду боту (старт, стоп...)
	return c.String(http.StatusOK, "ok\n")
}
func deletebot(c *echo.Context) error {
	//	TODO удалить бота
	return c.String(http.StatusOK, "ok\n")
}
