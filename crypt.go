package main


type MPPC struct {
	packedBytes		[]byte
	unpackedBytes	[]byte
	packedOffset	byte
	code1 			int
	code2 			int
	code3 			int
	code4 			int
}