package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type userLoginStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (S *HTTPServer) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("/login")
	var requestBody userLoginStruct
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := S.app.LoginUser(requestBody.Username, requestBody.Password)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(*user)
}
