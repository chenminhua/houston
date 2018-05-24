package utils

import (
	"net"
	"errors"
)

func GetFileDescriptor(conn net.Conn) (fd uintptr, err error) {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return 0, errors.New("not a TCPConn")
	}

	file, err := tcpConn.File()
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return file.Fd(), nil
}
