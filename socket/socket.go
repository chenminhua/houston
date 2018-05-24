package socket

import (
	"houston/netpoll"
	"fmt"
	"net/http"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"reflect"
	"net"
	"log"
)

// golang内建的map不是并发安全的
// 解决方案一     （csp）
// 解决方案二     （加锁）
// 解决方案三     （sync.Map）

var ConnMap = make(map[string]net.Conn)

func
SetUpWebSocketServer() {

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

		//fd, _:= getFileDescriptor(conn)
		//		//println(fd)
		//		//desc := netpoll.NewDesc(fd, netpoll.EventRead)

		desc, err := netpoll.Handle(conn, netpoll.EventRead | netpoll.EventWrite | netpoll.EventEdgeTriggered)
		HandleError(err)

		poller.Start(desc, func(ev netpoll.Event) {
			println("something")
			msg, op, err := wsutil.ReadClientData(conn)
			if reflect.TypeOf(err) == reflect.TypeOf(wsutil.ClosedError{}) {
				println("cleanUp", op, authToken)
				cleanUpConn(conn, authToken)
			} else {
				HandleError(err)
			}
			if op == ws.OpText {
				err = wsutil.WriteServerMessage(conn, ws.OpText, msg)
				HandleError(err)
			}
		})
	}))
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
