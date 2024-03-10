package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (S *HTTPServer) HandleStartGame(w http.ResponseWriter, r *http.Request) {
	log.Println("/start-game")
	username := mux.Vars(r)["username"]
	game, err := S.app.StartGame(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&game)
}
func (S *HTTPServer) HandleDrawCard(w http.ResponseWriter, r *http.Request) {
	log.Println("/draw")
	username := mux.Vars(r)["username"]
	game := S.app.UserDrawsCard(username)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(game)
}

func (S *HTTPServer) HandlerLeaderBoard(w http.ResponseWriter, r *http.Request) {
	user := S.app.GetLeaderBoard()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
