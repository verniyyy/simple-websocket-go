package src

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

// Serve ...
func Serve() {
	const (
		host = ""
		port = 3000
	)
	addr := fmt.Sprintf("%s:%d", host, port)
	r := NewRouter()

	log.Printf("listning on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Printf("Error: %v", err)
	}
}

// Router ...
type Router struct {
	*chi.Mux
}

// NewRouter ...
func NewRouter() Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ws", wsHandler)

	return Router{r}
}

// wsConn is function of upgrade to websocket connection from http connection
var wsConn = (&websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}).Upgrade

// wsHandler ...
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsConn(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Client Connected")

	err = conn.WriteMessage(websocket.TextMessage, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	reader(conn)
}

// reader ...
func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))

		// echo received message
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
