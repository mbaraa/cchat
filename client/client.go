package client

import (
	"bufio"
	"fmt"
	"os"

	ws "github.com/gorilla/websocket"
)

func Start(roomID, serverAddress string) {
	conn, _, err := ws.DefaultDialer.Dial(
		fmt.Sprintf("ws://%s/join-room/?room_id=%s", serverAddress, roomID), nil)

	if err != nil {
		panic("room not found or wrong server address")
	}

	defer conn.Close()

	fmt.Println("starting client, room id:", roomID)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("enter msg: ")
			in, _ := reader.ReadString('\n')
			conn.WriteMessage(ws.TextMessage, []byte(in))
		}
	}()

	for {
		_, msg, _ := conn.ReadMessage()

		fmt.Println("\nreceived:", string(msg))
		fmt.Print("enter msg: ")
	}
}
