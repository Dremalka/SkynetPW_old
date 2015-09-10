package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"fmt"
)

// api

// отправляет json-ответ с массивом текущих ботов
func listBot(c *echo.Context) error {
	//	TODO сформировать список ботов и отправить клиенту
	// временные данные
	list, err := MB.GetListBots()
	if err != nil {
		return c.JSON(http.StatusConflict, nil)
	}
	return c.JSON(http.StatusOK, list)
}

func createBot(c *echo.Context) error {
	//	TODO создать нового бота
	return c.String(http.StatusOK, "ok\n")
}

// sendActionToBot отправить команду боту
func sendActionToBot(c *echo.Context) error {
	param := make(map[string]interface{})
	id := c.Param("id")			// распарсить и получить id
	action := c.Param("action")	// и команду

	switch action {
	case "update":
		ProcessID, _ := strconv.Atoi(c.Query("ProcessID"))
		param["ProcessID"] = ProcessID
	case "connect":
		infbot := make(map[string]string)
		infbot["name"] = "bot1"
		id, _ = MB.AddBot(infbot)
	case "disconnect":
//		TODO отключить бота от игрового клиента
	}

	err := MB.SendActionToBot(id, action, param)	// отправить команду боту
	if err != nil {
		fmt.Println(err)
	}

	return c.String(http.StatusOK, "ok\n")
}
func deleteBot(c *echo.Context) error {
	//	TODO удалить бота
	return c.String(http.StatusOK, "ok\n")
}
