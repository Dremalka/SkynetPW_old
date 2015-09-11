package main
import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"fmt"
)

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

func randomNextBytes(count int) ([]byte, error) {
	b := make([]byte, count)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getRC4Key(encoderhash, uid, hash []byte) []byte {
	nhash := append(hash, encoderhash...)
	hmacmd5 := hmac.New(md5.New, uid)
	hmacmd5.Write(nhash)
	data := hmacmd5.Sum(nil)
	result := make([]byte, len(data))
	copy(result, data)
	return result
}

func newMPPC() *MPPC {
	mppc := &MPPC{}
	return mppc
}

// UnpackByte расшифровывает байт из очереди по модернизированному алгоритму MPPC
func (mppc *MPPC) UnpackByte(packedByte byte) []byte {
	code1 := mppc.code1
	code2 := mppc.code2
	code3 := mppc.code3
	code4 := mppc.code4

	mppc.packedBytes = append(mppc.packedBytes, packedByte)
	var unpackedChunk []byte
	var tempbuf []byte

	if len(mppc.unpackedBytes) >= 10240 {
		// удалить первые 2048 байт
		fmt.Println("Очистить первые 2048 байт.")
		tempbuf = make([]byte, len(mppc.unpackedBytes)-2048)
		copy(tempbuf, mppc.unpackedBytes[2048:])
		mppc.unpackedBytes = tempbuf
	}

	loop:
	for {
		switch code3 {
		case 0:
			if mppc.hasbits(4) == true {
				if mppc.getpackedbits(1) == 0 {
					code1 = 1
					code3 = 1
				} else {
					if mppc.getpackedbits(1) == 0 {
						code1 = 2
						code3 = 1
					} else {
						if mppc.getpackedbits(1) == 0 {
							code1 = 3
							code3 = 1
						} else {
							if mppc.getpackedbits(1) == 0 {
								code1 = 4
								code3 = 1
							} else {
								code1 = 5
								code3 = 1
							}
						}
					}
				}
			} else {
				break loop
			}
		case 1:
			switch code1 {
			case 1:
				if mppc.hasbits(7) == true {
					outB := byte(mppc.getpackedbits(7))
					unpackedChunk = append(unpackedChunk, outB)
					mppc.unpackedBytes = append(mppc.unpackedBytes, outB)
					code3 = 0
				} else {
					break loop
				}
			case 2:
				if mppc.hasbits(7) == true {
					outB := byte(mppc.getpackedbits(7) | 0x80)
					unpackedChunk = append(unpackedChunk, outB)
					mppc.unpackedBytes = append(mppc.unpackedBytes, outB)
					code3 = 0
				} else {
					break loop
				}
			case 3:
				if mppc.hasbits(13) == true {
					code4 = int(mppc.getpackedbits(13)) + 0x140
					code3 = 2
				} else {
					break loop
				}
			case 4:
				if mppc.hasbits(8) == true {
					code4 = int(mppc.getpackedbits(8)) + 0x40
					code3 = 2
				} else {
					break loop
				}
			case 5:
				if mppc.hasbits(6) == true {
					code4 = int(mppc.getpackedbits(6))
					code3 = 2
				} else {
					break loop
				}
			}
		case 2:
			if code4 == 0 {
				if mppc.packedOffset != 0 {
					mppc.packedOffset = 0
					// удалить первый байт в mppc.packedBytes
					tempbuf = make([]byte, len(mppc.packedBytes)-1)
					copy(tempbuf, mppc.packedBytes[1:])
					mppc.packedBytes = tempbuf
				}
				code3 = 0
				continue loop
			}
			code2 = 0
			code3 = 3
		case 3:
			if mppc.hasbits(1) == true {
				if mppc.getpackedbits(1) == 0 {
					code3 = 4
				} else {
					code2++
				}
			} else {
				break loop
			}
		case 4:
			var copySize int
			if code2 == 0 {
				copySize = 3
			} else {
				size := code2 + 1
				if mppc.hasbits(size) == true {
					copySize = int(mppc.getpackedbits(size)) + (1 << uint(size))
				} else {
					break loop
				}
			}
			unpackedChunk = mppc.CopyArray(code4, copySize, unpackedChunk)
			code3 = 0
		}
	}
	mppc.code1 = code1
	mppc.code2 = code2
	mppc.code3 = code3
	mppc.code4 = code4

	return unpackedChunk
}

// Unpack метод расшифровывает массив байтов по модернизированному алгоритму MPPC
func (mppc *MPPC) Unpack(compressedBytes []byte) []byte {
	var rtnList []byte
	for _, b := range compressedBytes {
		rtnList = append(rtnList, mppc.UnpackByte(b)...)
	}
	return rtnList
}

// CopyArray метод отсекает часть массива. Остаток переносит в начало массива
func (mppc *MPPC) CopyArray(shift, size int, unpackedChunkData []byte) []byte {
	for i := 0; i < size; i++ {
		pIndex := len(mppc.unpackedBytes) - shift
		if pIndex < 0 {
			return unpackedChunkData
		}
		b := mppc.unpackedBytes[pIndex]
		mppc.unpackedBytes = append(mppc.unpackedBytes, b)
		unpackedChunkData = append(unpackedChunkData, b)
	}
	return unpackedChunkData
}

func (mppc *MPPC) getpackedbits(bitCount int) uint {
	if bitCount > 16 {
		return 0
	}

	if mppc.hasbits(bitCount) == false {
		panic(fmt.Sprintln("Unpack bit stream overflow"))
	}

	alBitCount := bitCount + int(mppc.packedOffset)
	alByteCount := (alBitCount + 7) / 8

	var v uint32
	for i := 0; i < alByteCount; i++ {
		v |= uint32(mppc.packedBytes[i]) << uint32(24-i*8)
	}

	v <<= mppc.packedOffset
	v >>= uint32(32 - bitCount)

	mppc.packedOffset += byte(bitCount)
	freeBytes := mppc.packedOffset / 8

	if freeBytes != 0 {
		// удалить первые n-байт
		//fmt.Printf("getpackedbits. Удалить первые %d байт. mppc.packedBytes = [% X]\n", freeBytes, mppc.packedBytes)
		tempbuf := make([]byte, len(mppc.packedBytes)-int(freeBytes))
		copy(tempbuf, mppc.packedBytes[int(freeBytes):])
		mppc.packedBytes = tempbuf
		//fmt.Printf("getpackedbits. Новый mppc.packedBytes = [% X]\n", mppc.packedBytes)
	}
	mppc.packedOffset %= 8
	return uint(v)
}

func (mppc *MPPC) hasbits(count int) bool {
	if len(mppc.packedBytes)*8-int(mppc.packedOffset) >= count {
		return true
	}
	return false
}
