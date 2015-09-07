package main

import (
	"testing"
)

func Test_getlistbot(t *testing.T) {
	list, err := getlistbot()
	if err != nil {
		t.Error(err)
	}
	// результат не должен равняться nil
	if list == nil {
		t.Errorf("Массив list не инициализирован.\n")
	}

	// размер массива должен равняться 2
	lenorig := 2
	if len(list) != lenorig {
		t.Errorf("Размер массива равен %d. Должен быть равен %d\n", len(list), lenorig)
	} else {
		ID0 := "123"
		if list[0].ID != ID0 {
			t.Errorf("list[0].ID = %s. Должен быть - %s", list[0].ID, ID0)
		}
		Name0 := "abc"
		if list[0].Name != Name0 {
			t.Errorf("list[0].Name = %s. Должен быть - %s", list[0].Name, Name0)
		}
		ID1 := "456"
		if list[1].ID != ID1 {
			t.Errorf("list[1].ID = %s. Должен быть - %s", list[1].ID, ID1)
		}
		Name1 := "dfg"
		if list[1].Name != Name1 {
			t.Errorf("list[1].Name = %s. Должен быть - %s", list[1].Name, Name1)
		}

	}

}
