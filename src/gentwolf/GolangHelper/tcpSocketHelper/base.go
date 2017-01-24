package tcpSocketHelper

import (
	"bytes"
	"net"
)

const (
	NetTypeTCP   = "tcp"
	NetTypeTCP4  = "tcp4"
	NetTypeTCP6  = "tcp6"
	HeaderLength = 4
)

func pack(msg []byte) []byte {
	b := bytes.Buffer{}
	b.Write(uintToBytes(uint32(len(msg))))
	b.Write(msg)
	return b.Bytes()
}

func uintToBytes(i uint32) []byte {
	b := []byte{0, 0, 0, 0}
	b[0] = byte(i)
	b[1] = byte(i >> 8)
	b[2] = byte(i >> 16)
	b[3] = byte(i >> 24)

	return b
}

func bytesToUint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func readStream(conn *net.Conn, callback func(b []byte)) {
	buffer := bytes.Buffer{}
	dataLength := 0
	for {
		tmp := make([]byte, 1024)
		n, err := (*conn).Read(tmp)
		if err != nil {
			break
		}

		buffer.Write(tmp[0:n])

		for {
			msg := buffer.Bytes()
			msgLength := buffer.Len()
			if dataLength == 0 && msgLength >= HeaderLength {
				dataLength = int(bytesToUint32(msg[0:HeaderLength]))
				if msgLength-HeaderLength >= dataLength {
					offset := HeaderLength + dataLength
					callback(msg[HeaderLength:offset])

					buffer.Truncate(0)
					buffer.Write(msg[offset:])
					dataLength = 0
				}
			} else {
				break
			}
		}
	}
}
