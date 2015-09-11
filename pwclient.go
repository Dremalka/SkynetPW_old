package main
import (
	"net"
	"crypto/rc4"
	"fmt"
)

// PWClient объект игрового клиента низкого уровня
type PWClient struct {
	login			string
	password		string
	uid 			[]byte
	uid2			[]byte
	token 			[]byte
	conn 			net.Conn
	key 			[]byte
	isLoginCompleted	bool
	loginResult			bool
	rc4encode			[]byte
	rc4decode			[]byte
	rccipherenc			*rc4.Cipher
	rccipherdec			*rc4.Cipher
	mppc				*MPPC
	accountkey			[]byte
	unkIDOnlineAnnounce []byte
	unkDataOnlineAnnounce []byte
}

// newPWClient создает новый игровой клиент низкого уровня
func newPWClient() *PWClient {
	pw := &PWClient{}
	return pw
}

// Connect метод подключается к игровому серверу
func (pw *PWClient) Connect(server, login, pass string) error {
	pw.login = login
	pw.password = pass

	// отправить запрос серверу mail.ru и получить уиды и токен
	uid,uid2,token,err := getLoginPass(login, pass)
	if err != nil {
		return err
	}

	pw.uid = []byte(uid)
	pw.uid2 = []byte(uid2)
	pw.token = []byte(token)

	connectViaSocks := 0	// флаг определяющий соединенеие с игровым сервером (напрямую, через прокси-сокет)
	var conn net.Conn
	if connectViaSocks == 1 {
		fmt.Println("Connect to socks.")
		conn, err = dialSocks4(SOCKS4, "localhost:30000", server)
		if err != nil {
			fmt.Println("socks error:", err)
			return err
		}
	} else {
		conn, err = net.Dial("tcp", server)
		if err != nil {
			fmt.Println("dial error:", err)
			return err
		}
	}

	pw.conn = conn

	go getFromServer(pw)

	return nil

}