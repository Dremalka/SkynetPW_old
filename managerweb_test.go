package main

import(
	"testing"
)

func Test_newManagerWeb(t *testing.T) {
	mw, err := newManagerWeb()
	if err != nil {
		t.Error("Типовая ошибка при создании нового объекта.")
	}
	if mw == nil {
		t.Error("метод newManagerWeb() не должен возвращать nil.")
	}
}
