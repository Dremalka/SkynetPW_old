package main

import(
	"testing"
	//"net/http"
)

func Test_newManagerWeb(t *testing.T) {
	mw, err := newManagerWeb()
	if err != nil {
		t.Errorf("Типовая ошибка при создании нового объекта. err = %v\n", err)
	}
	if mw == nil {
		t.Error("метод newManagerWeb() не должен возвращать nil.")
	}
}

//func Test_Start(t *testing.T) {
//	// создать менеджер веб-интерфейса
//	mw, err := newManagerWeb()
//	if err != nil || mw == nil {
//		t.Errorf("Невозможно проверить создание сервера, т.к. невозможно создать объект Менеджер веб-интерфейса. err = %v\n", err)
//	}
//	mw.Start()	// запустить сервер
//
//	// запросить главную страницу
//	request, err := http.NewRequest("GET", "http://localhost:8080", nil)
//	res, err := http.DefaultClient.Do(request)
//	if err != nil {
//		t.Error(err)	// произошла ошибка пока отправлялся запрос
//	}
//
//	// статус ответа должен быть 200
//	if res.StatusCode != 200 {
//		t.Errorf("Получен код: %d. Ожидается: 200.\n", res.StatusCode)
//	}
//}
