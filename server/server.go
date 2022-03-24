package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

var (
	upgrader = ws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Room struct {
	roomID string
	conns  map[string]*ws.Conn
}

var (
	rooms = map[string]Room{}
)

func listenToMsgs(conn *ws.Conn, room Room) {
	for {
		mType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(room.conns, conn.RemoteAddr().String())
			break
		}
		log.Printf("Received: %s, from: %s", msg, conn.RemoteAddr().String())
		for _, _conn := range room.conns {
			if _conn.RemoteAddr().String() != conn.RemoteAddr().String() {
				err = _conn.WriteMessage(mType, msg)
			}
			if err != nil {
				log.Println("error writing message", err)
				break
			}
		}
	}
}

func joinRoom(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("room_id")
	room, exists := rooms[id]
	if exists {
		conn, err := upgrader.Upgrade(res, req, nil)
		if err != nil {
			log.Println("error reading message", err)
			return
		}

		room.conns[conn.RemoteAddr().String()] = conn
		listenToMsgs(conn, room)
	}

}

func createRoom(res http.ResponseWriter, req *http.Request) {
	id := uuid.NewString()[:8]
	newRoom := Room{id, make(map[string]*ws.Conn)}

	rooms[id] = newRoom

	json.NewEncoder(res).Encode(map[string]string{
		"room_id": id,
	})
}

func Start(port string) {
	http.HandleFunc("/create-room/", createRoom)
	http.HandleFunc("/join-room/", joinRoom)

	fmt.Println("starting server on port:", port)
	http.ListenAndServe(":"+port, nil)
}
