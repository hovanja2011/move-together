package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/hovanja2011/move-together/internal/config"
	"github.com/hovanja2011/move-together/internal/user"
	"github.com/hovanja2011/move-together/internal/user/db"
	"github.com/hovanja2011/move-together/pkg/client/mongodb"
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

	cfg := config.GetConfig()
	cfgMongo := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username,
		cfgMongo.Password, cfgMongo.Database, cfgMongo.AuthDB)
	if err != nil {
		panic(err)
	}

	storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	users, err := storage.FindAll(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	user1 := user.User{
		ID:           "",
		Email:        "user1@gmail.com",
		Username:     "user1name",
		PasswordHash: "12345",
	}
	user1ID, err := storage.Create(context.Background(), user1)
	if err != nil {
		panic(err)
	}
	logger.Info(user1ID)

	user2 := user.User{
		ID:           "",
		Email:        "user2@gmail.com",
		Username:     "user2name",
		PasswordHash: "67890",
	}
	user2ID, err := storage.Create(context.Background(), user2)
	if err != nil {
		panic(err)
	}
	logger.Info(user2ID)

	user2Found, err := storage.FindOne(context.Background(), user2ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(user2Found)

	user2Found.Email = "newEmail@mail.ru"
	err = storage.Update(context.Background(), user2Found)
	if err != nil {
		panic(err)
	}

	err = storage.Delete(context.Background(), user2ID)
	if err != nil {
		panic(err)
	}
	// _, err = storage.FindOne(context.Background(), user2ID)
	// if err != nil {
	// 	panic(err)
	// }

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	logger.Info("start router")
	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenError error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenError = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket %s", socketPath)

	} else {
		logger.Info("listen tcp")
		listener, listenError = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenError != nil {
		logger.Fatal(listenError)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 26 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}
