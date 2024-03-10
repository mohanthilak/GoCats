package main

import (
	redispkg "GoCats/Redis"
	runner "GoCats/Runner"
	user "GoCats/User"
	"GoCats/server"
	"flag"
	"log"
	"os"
	"strconv"
)

func main() {
	wsChan := make(chan []map[string]interface{})
	defaultPort := 8000
	port := flag.Int("port", defaultPort, "port")

	pp := os.Getenv("PORT")
	var pp2 = *port
	var err error
	if pp != "" {
		pp2, err = strconv.Atoi(pp)
		if err != nil {
			log.Println(err)
		}
	}

	port = &pp2

	redisController := redispkg.NewRedisClient()

	userModule := user.CreateUserModule(redisController)

	app := runner.NewRunner(*userModule, wsChan)

	httpServer := server.NewHTTPServer(*port, app, wsChan)
	httpServer.StartServer()

}
