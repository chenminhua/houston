package restapi

import (
	"net/http"
	"log"
	"fmt"
	"github.com/gorilla/mux"
)

type Server struct {
	Port string
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

	fmt.Fprintln(w, "pub")
}
