package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type ConnectUser struct {
	Websocket *websocket.Conn
	ClientIP  string
}

func newConnectUser(ws *websocket.Conn, clientIP string) *ConnectUser {
	return &ConnectUser{
		Websocket: ws,
		ClientIP:  clientIP,
	}
}

var users = make(map[ConnectUser]int)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	//t, _ := template.ParseFiles("index.html")
	//p, _ := loadChat()
	//t.Execute(w, p)
	http.ServeFile(w, r, "index.html")
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	var socketClient *ConnectUser = newConnectUser(conn, conn.RemoteAddr().String())
	users[*socketClient] = 0
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			delete(users, *socketClient)
			return
		}

		for client := range users {
			client.Websocket.WriteMessage(messageType, message)
		}

	}
}

func main() {
	http.HandleFunc("/chat/", chatHandler)
	http.HandleFunc("/ws", ws)
	http.ListenAndServe(":8080", nil)
}
