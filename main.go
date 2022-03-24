package main

import (
	"cchat/client"
	"cchat/server"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	mode := flag.String("m", "client", "Run in server, or client mode")

	// client stuff
	serverAddress := flag.String("s", "localhost:8080", "Chat server")
	roomID := flag.String("r", "0", "Room ID, must be provided!")

	// server stuff
	port := flag.Int("p", 8080, "Run server on the given port, default is 8080")

	flag.Parse()
	switch *mode {
	case "client", "":
		if *roomID == "0" {
			resp, _ := http.Get("http://" + *serverAddress + "/create-room/")

			var respBody map[string]string
			json.NewDecoder(resp.Body).Decode(&respBody)
			*roomID = respBody["room_id"]

			fmt.Println("room id:", respBody["room_id"])
		}

		client.Start(*roomID, *serverAddress)

	case "server":
		if *port > 20 {
			server.Start(fmt.Sprint(*port))
		}
	}
}
