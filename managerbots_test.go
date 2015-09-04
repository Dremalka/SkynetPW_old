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
			t.Errorf("Не инициализирован массив списка ботов. %q != %q", mb.ListBot, someval)
		}
	}


}
