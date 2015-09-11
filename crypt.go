package main

// MPPC структура объекта для распаковки пакета данных от игрового сервера
type MPPC struct {
	packedBytes		[]byte
	unpackedBytes	[]byte
	packedOffset	byte
	code1 			int
	code2 			int
	code3 			int
	code4 			int
}