package main
import (
	"net"
	"strconv"
	"crypto/rc4"
	"fmt"
	"strings"
	"errors"
	"net/url"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/xml"
	"encoding/binary"
	"io"
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

//	TODO удалить
	// временно. для отладки
	rows, err := MB.db.Query("SELECT server, login, password FROM servers")
	rows.Next()
	err = rows.Scan(&server, &login, &pass)


	pw.login = login
	pw.password = pass

	// отправить запрос серверу mail.ru и получить уиды и токен
	uid,uid2,token,err := getLoginAndPass(login, pass)
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

// getLoginAndPass функция отправляет запрос на сайт mail.ru и получает уиды и токен для авторизации на игровом сайте
func getLoginAndPass(email, pass string) (string, string, string, error) {
	split := strings.Split(email, "@")
	if len(split) < 2 {
		return "", "", "", errors.New("Неправильный e-mail.")
	}

	domain := split[1]

	mailDomains := []string{"mail.ru", "inbox.ru", "bk.ru", "list.ru"}
	foundMailDomains := false
	for _, v := range mailDomains {
		if v == domain {
			foundMailDomains = true
			break
		}
	}

	var uid, uid2, token string

	if foundMailDomains == true {
		fmt.Println("Found mail domain")
//		TODO обработка ситуации, когда e-mail заведен в mail.ru
	} else {
		apiURL := "http://authdl.mail.ru"
		resource := "/sz.php"
		data := url.Values{}
		params := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><AutoLogin ProjectId="61" SubProjectId="0" ShardId="0" Username="%s" Password="%s"/>`, email, pass)
		u, _ := url.ParseRequestURI(apiURL)
		u.Path = resource
		urlStr := fmt.Sprintf("%v", u)

		client := &http.Client{}
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(params)) // <-- URL-encoded payload
		r.Header.Add("User-Agent", "Downloader/4260")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		r.Header.Add("Accept-Encoding", "identity")
		resp, _ := client.Do(r)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", "", "", err
		}

		type AuthPers struct {
			XMLName xml.Name `xml:"AutoLogin"`
			UID2    string   `xml:"PersId,attr"`
			Token   string   `xml:"Key,attr"`
		}
		var q AuthPers
		xml.Unmarshal(body, &q)
		uid2 = q.UID2
		token = q.Token

		//
		params = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><PersList ProjectId="61" SubProjectId="0" ShardId="0" Username="%s" Password="%s"/>`, email, pass)
		r, _ = http.NewRequest("POST", urlStr, bytes.NewBufferString(params)) // <-- URL-encoded payload
		r.Header.Add("User-Agent", "Downloader/4260")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		r.Header.Add("Accept-Encoding", "identity")
		resp, _ = client.Do(r)

		body, err = ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))

		type Pers struct {
			XMLName xml.Name `xml:"Pers"`
			ID      string   `xml:"Id,attr"`
			Title   string   `xml:"Title,attr"`
			Cli     string   `xml:"Cli,attr"`
		}

		type PersList struct {
			XMLName  xml.Name `xml:"PersList"`
			PersID   string   `xml:"PersId,attr"`
			PersList []Pers   `xml:"Pers"`
		}

		var q1 PersList
		xml.Unmarshal(body, &q1)
		if len(q1.PersList) == 0 {
			return "", "", "", errors.New("У учетной записи нет игровых аккаунтов.")
		}
		uid = q1.PersList[0].ID

	}

	//fmt.Println("Ответ от сервера", uid, uid2, token)
	return uid, uid2, token, nil
}

// getFromServer Обработчик данных от сервера. Получает данные и вызывает обработчики пакетов
func getFromServer(pwclient *PWClient) {
	// если выходим, то закрываем соединение
	defer pwclient.conn.Close()

	conn := pwclient.conn
	inBuf := make([]byte, 1024)
	for {

		n, err := conn.Read(inBuf[0:])
		if err != nil && err != io.EOF {
			panic(fmt.Sprintf("Ошибка при получении данных из соединения. %v", err))
		}
		if err == io.EOF {
			panic(fmt.Sprintf("Ошибка при получении данных из соединения. %v. bufer = %x", err, inBuf))
		}
		if n < 2 {
			panic(fmt.Sprintf("Пакет 0-ой длины."))
		}
		buf := make([]byte, n)
		copy(buf, inBuf)

		if pwclient.isLoginCompleted {
			fmt.Printf("До дешифрования (получение). Длина буфера = %d, [% X]\n", len(buf), buf)
			bufsrc := make([]byte, len(buf))
			copy(bufsrc, buf)
			pwclient.rccipherdec.XORKeyStream(buf, bufsrc)
			fmt.Printf("После дешифрования (получение). Длина буфера = %d, [% X]\n", len(buf), buf)
			bufsrc = make([]byte, len(buf))
			copy(bufsrc, buf)
			buf = pwclient.mppc.Unpack(bufsrc)
			fmt.Printf("После распаковки. Длина буфера = %d, [% X]\n", len(buf), buf)

		}

		reader := bytes.NewBuffer(buf)
		var firstbyte, secondbyte byte
		var kod, length uint16

		firstbyte, _ = reader.ReadByte()
		fmt.Printf("firstbyte (1) = %x\n", firstbyte)
		var src []byte
		if firstbyte >= 0x80 {
			firstbyte -= 0x80
			fmt.Printf("firstbyte (2) = %x\n", firstbyte)
			secondbyte, _ = reader.ReadByte()
			fmt.Printf("secondbyte = %x\n", secondbyte)
			src = append(src, firstbyte, secondbyte)
			kod = binary.BigEndian.Uint16(src)
		} else {
			kod = uint16(firstbyte)
		}

		fmt.Printf("kod = %x\n", kod)

		//kod := buf[0] // код пакета
		//var length uint16

		firstbyte = 0
		secondbyte = 0

		src = make([]byte, 0)
		firstbyte, _ = reader.ReadByte()
		if firstbyte >= 0x80 {
			firstbyte -= 0x80
			secondbyte, _ = reader.ReadByte()
			src = append(src, firstbyte, secondbyte)
			length = binary.BigEndian.Uint16(src)
		} else {
			length = uint16(firstbyte)
		}

		//if kod == 1 || kod == 2 || kod == 4 {
		//	length = int(buf[1]) + 2 // длина пакета

		if int(length) < reader.Len() {
			panic(fmt.Sprintf("Значение длины пакета больше, чем фактическая длина пакета."))
		}
		//}

		switch kod {
		case 1: // 0x01 ServerInfo
			fmt.Println("Код 1. ServerInfo", buf)
			packetServerInfo(pwclient, reader.Bytes())

		case 2:
			fmt.Println("Код 2. SMKey", buf)
			packetSMKey(pwclient, reader.Bytes())
		case 4:
			fmt.Println("Код 4. OnlineAnnounce", buf)
			packetOnlineAnnounce(pwclient, reader.Bytes())
		case 0x8f:
			fmt.Println("Код 0x8f. LastLogin", buf)
		case 0x53:
			fmt.Println("Код 0x53. LastLogin", buf)
			//packetRoleListRe(pwclient, reader.Bytes())
		default:
			fmt.Println("default", kod, reader.Len(), reader.Bytes())
		}

	}
}

func sendToServer(pwclient *PWClient, pkt []byte) error {
	conn := pwclient.conn
	fmt.Printf("До шифрования (отправка). [% X]\n", pkt)
	if pwclient.isLoginCompleted {
		pktsrc := pkt
		pwclient.rccipherenc.XORKeyStream(pkt, pktsrc)
		fmt.Printf("После шифрования (отправка). [% X]\n", pkt)
	}
	_, err := conn.Write(pkt)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// packetServerInfo обработчик пакета ServerInfo
func packetServerInfo(pwclient *PWClient, pkt []byte) {
	if pwclient.key != nil {
		return
	}
	keylen := pkt[0]
	key := pkt[1:keylen]
	pwclient.key = key
	fmt.Printf("%v, %v, % X\n", pkt, key, key)
	sendLogginAnnounce(pwclient)
}

// sendLoggingAnnounce Подготовка пакета LoggingAnnounce для отправки серверу
func sendLogginAnnounce(pwclient *PWClient) {
	hash := pwclient.token

	//	hmacsha256 := hmac.New(sha256.New, pwclient.key)
	//	hmacsha256.Write(append(pwclient.uid, pwclient.uid2...))
	//	result1 := hmacsha256.Sum(nil)
	//	result2 := hex.EncodeToString(result1)
	//	result3 := []byte(result2)
	//	hash := result3
	//	fmt.Printf("[% X]\n", pwclient.token)
	//	fmt.Printf("[% X]\n", hash)

	uid := pwclient.uid

	var data []byte
	data = append(data, uint8ToSlicebyte(uint8(len(uid)))...)
	data = append(data, []byte(uid)...)
	data = append(data, uint8ToSlicebyte(uint8(len(hash)))...)
	data = append(data, hash...)
	data = append(data, 0x02, 0x04, 0xFF, 0xFF, 0xFF, 0xFF)

	var pkt []byte
	pkt = append(pkt, byte(3))
	pkt = append(pkt, uint8ToSlicebyte(uint8(len(data)))...)
	pkt = append(pkt, data...)

	err := sendToServer(pwclient, pkt)
	if err != nil {

	}
}

func packetSMKey(pwclient *PWClient, pkt []byte) {
	fmt.Println(pkt)
	encHashLen := pkt[0]
	encHash := pkt[1 : encHashLen+1]
	force := pkt[encHashLen+1]
	fmt.Printf("%v\n", pkt)
	fmt.Printf("%v, %v, %v\n", encHashLen, encHash, force)

	pwclient.isLoginCompleted = true
	pwclient.loginResult = true

	// дек хэш
	decHash, err := randomNextBytes(16)
	//decHash := []byte{0x4D, 0x64, 0x6E, 0xD2, 0xCF, 0x4B, 0xB1, 0x5B, 0x7B, 0xB3, 0x70, 0x7B, 0x46, 0x10, 0x5C, 0xEB}
	//err := error(nil)
	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	rc4encode := getRC4Key(encHash, pwclient.uid, pwclient.token)
	rc4decode := getRC4Key(decHash, pwclient.uid, pwclient.token)

	pwclient.rc4encode = rc4encode
	pwclient.rc4decode = rc4decode

	rccipher, _ := rc4.NewCipher(rc4encode)
	pwclient.rccipherenc = rccipher
	rccipher, _ = rc4.NewCipher(rc4decode)
	pwclient.rccipherdec = rccipher

	mppc := newMPPC()
	pwclient.mppc = mppc

	//fmt.Printf("packetSMKey: [% X], [% X], [% X], [% X]\n", rc4encode, rc4decode, encode.m_Table, decode.m_Table)
	err = sendCMKey(pwclient, decHash)
	if err != nil {
		panic(fmt.Sprintf("Ошибка при отправке пакета CMKey"))
	}
}

func sendCMKey(pwclient *PWClient, decHash []byte) error {
	force := 1
	var data []byte
	data = append(data, uint8ToSlicebyte(uint8(len(decHash)))...)
	data = append(data, decHash...)
	data = append(data, byte(force))

	var pkt []byte
	pkt = append(pkt, byte(2))
	pkt = append(pkt, uint8ToSlicebyte(uint8(len(data)))...)
	pkt = append(pkt, data...)

	err := sendToServer(pwclient, pkt)
	if err != nil {

	}
	return nil
}

func packetOnlineAnnounce(pwclient *PWClient, pkt []byte) {
	accountkey := pkt[:4]
	//unkID := pkt[4:8]
	//unkData := pkt[8:]

	pwclient.accountkey = accountkey

	err := sendRoleList(pwclient, []byte{0xFF, 0xFF, 0xFF, 0xFF})
	if err != nil {
		panic(fmt.Sprintf("Ошибка при отправке пакета RoleList"))
	}
}

func sendRoleList(pwclient *PWClient, slot []byte) error {
	var data []byte
	data = append(data, pwclient.accountkey...)
	data = append(data, []byte{0, 0, 0, 0}...)
	data = append(data, slot...)

	var pkt []byte
	pkt = append(pkt, byte(0x52))
	pkt = append(pkt, uint8ToSlicebyte(uint8(len(data)))...)
	pkt = append(pkt, data...)

	err := sendToServer(pwclient, pkt)
	if err != nil {

	}
	return nil
}

func uint8ToSlicebyte(num uint8) []byte {

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		return nil
	}
	result := buf.Bytes()
	return result[:]
}
