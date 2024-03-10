package main

import (
	redispkg "GoCats/Redis"
	runner "GoCats/Runner"
	user "GoCats/User"
	"GoCats/server"
	"flag"
)

func main() {
	wsChan := make(chan []map[string]interface{})
	defaultPort := 8000
	port := flag.Int("port", defaultPort, "port")

	redisController := redispkg.NewRedisClient()

	userModule := user.CreateUserModule(redisController)

	app := runner.NewRunner(*userModule, wsChan)

	httpServer := server.NewHTTPServer(*port, app, wsChan)
	httpServer.StartServer()

}
