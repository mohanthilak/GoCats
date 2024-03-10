package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (S *HTTPServer) HandleSocketRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")
	username := mux.Vars(r)["username"]

	// Currently it's set to all origin, need to fetch the domain from the request and move forward accordingly
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade connection to full duplex")
	}

	// Currently uid is hardcoded to "123", need to fetch the userID from the requeset.
	client := newClient(conn, S.manager, username)
	S.manager.AddClient(client)

	log.Println("\n NEW CLIENT", S.manager)
	go client.ReadMessage()
	go client.WriteMessage()
}

func (S *HTTPServer) CheckAvailable(clinetID string) bool {
	if _, ok := S.manager.Clients[clinetID]; ok {
		return true
	}
	return false
}

func (S *HTTPServer) SendMessage(clientID, msgType string, mes []byte) error {
	if isOnline := S.CheckAvailable(clientID); !isOnline {
		log.Printf("User: %s Offline", clientID)
		return errors.New("user offline")
	}
	event := Event{
		Type:    msgType,
		Payload: mes,
	}
	log.Printf("%+v", S.manager)
	S.manager.Clients["123"].engres <- event
	return nil
}
func (S *HTTPServer) BroadCastMessage(msgType string, mes any) error {
	event := Event{
		Type:    msgType,
		Payload: mes,
	}
	log.Printf("%+v", S.manager)
	gg := S.manager.Clients
	for k, v := range gg {
		log.Println("user", k)
		v.engres <- event
	}
	return nil
}
