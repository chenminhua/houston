package socket

import (
	"houston/netpoll"
	"fmt"
	"net/http"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"reflect"
	"net"
	"errors"
	"log"
)


var ConnMap = make(map[string]net.Conn)

func SetUpWebSocketServer() {

	poller, err := netpoll.New(nil)
	if err != nil {
		fmt.Println(err)
	}

	http.ListenAndServe(":8087", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 在UpGrade之前完成鉴权
		authToken := r.URL.Query().Get("auth")
		conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
		ConnMap[authToken] = conn
		HandleError(err)

		fd, _:= getFileDescriptor(conn)
		desc := netpoll.NewDesc(fd, netpoll.EventRead)

		poller.Start(desc, func(ev netpoll.Event) {
			msg, op, err := wsutil.ReadClientData(conn)
			if reflect.TypeOf(err) == reflect.TypeOf(wsutil.ClosedError{}) {
				cleanUpConn(conn, authToken)
			} else {
				HandleError(err)
			}
			if op == ws.OpText {
				err = wsutil.WriteServerMessage(conn, ws.OpText, msg)
				HandleError(err)
			}
		})

		// 给每个长连接创建一个goroutine (仅仅为了监视websocket的状态？？)
		// go receiveMessage(conn, authToken)
	}))
}


func getFileDescriptor(conn net.Conn) (fd uintptr, err error) {
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


func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func cleanUpConn(conn net.Conn, userId string) {
	delete(ConnMap, userId)
	conn.Close()
}

func receiveMessage(conn net.Conn, userId string) {
	defer cleanUpConn(conn, userId)
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if reflect.TypeOf(err) == reflect.TypeOf(wsutil.ClosedError{}) {
			break;
		} else {
			HandleError(err)
		}
		if op == ws.OpText {
			err = wsutil.WriteServerMessage(conn, ws.OpText, msg)
			HandleError(err)
		}
	}
}
