package main

import (
	"fmt"
	"go/src/log"
)

//MB глобальная переменная, которая содержит объект Менеджер ботов
var MB *ManagerBots

func main() {
	var err error

	MB, err = newManagerBots()
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintln(MB)

	var response string
	fmt.Println("Press Enter")
	_, _ = fmt.Scanln(&response)
	fmt.Println("Exit.")
}
