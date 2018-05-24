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
	"io"
)

var CPoll = ConnPool{
	ConnMap:make(map[string]net.Conn),
}

func SetUpWebSocketServer() {

	poller, err := netpoll.New(nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("poller", poller)

	http.ListenAndServe(":8087", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 在UpGrade之前完成鉴权
		authToken := r.URL.Query().Get("auth")
		conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
		CPoll.Set(authToken, conn)
		HandleError(err)

		desc, err := netpoll.Handle(conn, netpoll.EventRead | netpoll.EventEdgeTriggered)
		HandleError(err)
		fmt.Println("desc", desc)

		poller.Start(desc, func(ev netpoll.Event) {
			println("something")
			msg, op, err := wsutil.ReadClientData(conn)
			fmt.Println(msg)
			fmt.Println(op)
			if (reflect.TypeOf(err) == reflect.TypeOf(wsutil.ClosedError{}) || err == io.EOF){
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
	CPoll.Delete(userId)
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
