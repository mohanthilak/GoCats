package server

import (
	runner "GoCats/Runner"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	router   *mux.Router
	port     int
	app      *runner.RunnerStruct
	manager  *manager
	messages <-chan []map[string]interface{}
}

func NewHTTPServer(port int, app *runner.RunnerStruct, msgChan <-chan []map[string]interface{}) *HTTPServer {
	return &HTTPServer{router: mux.NewRouter(), port: port, app: app, manager: MakeManager(), messages: msgChan}
}

func (S *HTTPServer) initiateRoutes() {
	S.router.HandleFunc("/user/login", S.HandleUserLogin)

	S.router.HandleFunc("/game/start/{username}", S.HandleStartGame)
	S.router.HandleFunc("/game/draw-card/{username}", S.HandleDrawCard)
	S.router.HandleFunc("/game/leader-board", S.HandlerLeaderBoard)

	S.router.HandleFunc("/ws/{username}", S.HandleSocketRoute)
	go func() {
		for {
			select {
			case msg := <-S.messages:
				fmt.Println("Received:", msg)
				// Process the received message
				S.BroadCastMessage("leaderBoard", msg)
			}
		}
	}()
}

// type routeHandlerType func(w http.ResponseWriter, r *http.Request)

//	func (S *HTTPServer) CreateHandler(fn routeHandlerType){
//		fn()
//	}
func (S *HTTPServer) StartServer() {
	S.initiateRoutes()

	server := http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(S.port),
		Handler:      handlers.CORS(handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}), handlers.AllowCredentials(), handlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "X-Requested-With", "Authorization", "Content-Type"}), handlers.AllowedOrigins([]string{"http://localhost:5173", "https://golden-griffin-81a6ad.netlify.app"}))(S.router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("server starting at port: ", S.port)
		err := server.ListenAndServe()
		if err != nil {
			log.Panicf("error while serving the server. %s", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c

	log.Printf("Gracefully Ending the Server after receiving signal %s", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
