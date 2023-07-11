package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hovanja2011/move-together/internal/user"
	"github.com/hovanja2011/move-together/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	w.Write([]byte(fmt.Sprintf(("hello %s"), name)))
}

func main() {
	logger := logging.GetLogger()

	logger.Info("create router")
	router := httprouter.New()

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	logger.Info("start router")
	start(router)
}

func start(router *httprouter.Router) {
	logger := logging.GetLogger()
	logger.Info("start application")

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 26 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("server is listening port :1234")
	log.Fatalln(server.Serve(listener))
}
