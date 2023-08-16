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
	driver2 "github.com/hovanja2011/move-together/internal/driver"
	driver "github.com/hovanja2011/move-together/internal/driver/db"
	passenger2 "github.com/hovanja2011/move-together/internal/passenger"
	passenger "github.com/hovanja2011/move-together/internal/passenger/db"
	ride2 "github.com/hovanja2011/move-together/internal/ride"
	ride "github.com/hovanja2011/move-together/internal/ride/db"
	"github.com/hovanja2011/move-together/internal/user"

	"github.com/hovanja2011/move-together/pkg/client/postgresql"
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

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	repositoryRide := ride.NewRepository(postgreSQLClient, logger)
	repositoryPassenger := passenger.NewRepository(postgreSQLClient, logger)
	repositoryDriver := driver.NewRepository(postgreSQLClient, logger)

	logger.Info("register driver handler")
	driverHandler := driver2.NewHandler(repositoryDriver, logger)
	driverHandler.Register(router)

	logger.Info("register passenger handler")
	passengerHandler := passenger2.NewHandler(repositoryPassenger, logger)
	passengerHandler.Register(router)

	logger.Info("register ride handler")
	rideHandler := ride2.NewHandler(repositoryRide, logger)
	rideHandler.Register(router)

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
