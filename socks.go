package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

// варианты
const (
	SOCKS4 = iota
	SOCKS4A
	SOCKS5
)

func dialSocks4(socksType int, proxy, targetAddr string) (conn net.Conn, err error) {
	// dial TCP
	conn, err = net.Dial("tcp", proxy)
	if err != nil {
		return
	}

	// connection request
	host, port, err := splitHostPort(targetAddr)
	if err != nil {
		return
	}
	ip := net.IPv4(0, 0, 0, 1).To4()
	if socksType == SOCKS4 {
		ip, err = lookupIP(host)
		if err != nil {
			return
		}
	}
	req := []byte{
		4,                          // version number
		1,                          // command CONNECT
		byte(port >> 8),            // higher byte of destination port
		byte(port),                 // lower byte of destination port (big endian)
		ip[0], ip[1], ip[2], ip[3], // special invalid IP address to indicate the host name is provided
		0, // user id is empty, anonymous proxy only
	}
	if socksType == SOCKS4A {
		req = append(req, []byte(host+"\x00")...)
	}

	resp, err := sendReceive(conn, req)
	if err != nil {
		return
	} else if len(resp) != 8 {
		err = errors.New("Server does not respond properly.")
	}
	switch resp[1] {
	case 90:
	// request granted
	case 91:
		err = errors.New("Socks connection request rejected or failed.")
	case 92:
		err = errors.New("Socks connection request rejected becasue SOCKS server cannot connect to identd on the client.")
	case 93:
		err = errors.New("Socks connection request rejected because the client program and identd report different user-ids.")
	default:
		err = errors.New("Socks connection request failed, unknown error.")
	}
	return
}

func sendReceive(conn net.Conn, req []byte) (resp []byte, err error) {
	_, err = conn.Write(req)
	if err != nil {
		return
	}
	resp, err = readAll(conn)
	return
}

func readAll(conn net.Conn) (resp []byte, err error) {
	resp = make([]byte, 1024)
	n, err := conn.Read(resp)
	resp = resp[:n]
	return
}

func lookupIP(host string) (ip net.IP, err error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return
	}
	if len(ips) == 0 {
		err = fmt.Errorf("Cannot resolve host: %s.", host)
		return
	}
	ip = ips[0].To4()
	if len(ip) != net.IPv4len {
		fmt.Println(len(ip), ip)
		err = errors.New("IPv6 is not supported by SOCKS4")
		return
	}
	return
}

func splitHostPort(addr string) (host string, port uint16, err error) {
	host, portStr, err := net.SplitHostPort(addr)
	portInt, err := strconv.ParseUint(portStr, 10, 16)
	port = uint16(portInt)
	return
}
