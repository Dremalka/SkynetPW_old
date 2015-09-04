package main

import (
	"fmt"
	"go/src/log"
)

//MB глобальная переменная, которая содержит объект Менеджер ботов
var MB *ManagerBots
//MW глобальная переменная, которая содержит объект Менеджер веб-интерфейса
var MW *ManagerWeb

func main() {
	var err error

	MB, err = newManagerBots()
	if err != nil {
		log.Println(err)
	}

	MW, err = newManagerWeb()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(MB, MW)

	var response string
	fmt.Println("Press Enter")
	_, _ = fmt.Scanln(&response)
	fmt.Println("Exit.")
}
