package restapi

import (
	"net/http"
	"log"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"io"
	"encoding/json"
	"houston/socket"
	"strconv"
	"github.com/gobwas/ws/wsutil"
	"github.com/gobwas/ws"
)

type Server struct {
	Port string
}

type Message struct {
	Content string `json:"content"`
	Receivers []int `json:"receivers"`
}

func (s *Server) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/ping", PingHandler)
	router.HandleFunc("/pub", PubHandler)
	http.Handle("/", router)
	println(s.Port)
	log.Fatal(http.ListenAndServe(s.Port, nil))
}

func PingHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong");
}

func PubHandler (w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	var message Message
	if err := json.Unmarshal(body, &message); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	for _, i := range message.Receivers {
		println(i)
		conn := socket.ConnMap[strconv.Itoa(i)]
		if conn != nil {
			err = wsutil.WriteServerMessage(conn, ws.OpText, []byte(message.Content))
			println("err...", err)
		}
	}

	fmt.Fprintln(w, "pub")
}
