package main
import (
	"net"
	"crypto/rc4"
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