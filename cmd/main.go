package main

import (
	"L0"
	"L0/nats"
	handler "L0/pkg/handler"
	"L0/pkg/repository"
	"L0/pkg/service"
	"context"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	//"github.com/spf13/viper"
	"github.com/subosito/gotenv" 
)

func main() {

	 if err := gotenv.Load(); err != nil {
	 	log.Fatalf("failed loading env variables: %s", err.Error())
	 }

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   "postgres",
		SSLMode:  "disable",
	})

	if err != nil {
		fmt.Println("Error while connecting to DB:", err.Error())
		return
	}

	repos := repository.NewRepository(db)
	service, err := service.NewService(repos)
	if err != nil {
		logrus.Errorf("Error while creating service: %s", err.Error())
	}
	handler := handler.NewHandler(service)

	srv := new(L0.Server)
	
	go func(){
		if err := srv.Run("8080", handler.InitRoutes()); err != nil {
			log.Fatalf("error occured while runnig http server")
		}
	}()

	clusterId   := "test-cluster"
	clientId    := "vibe"
	channelName := "vibeChannel"

	nats, err := nats.NewSubscribeToChannel(clusterId, clientId, channelName, repos, service)

	if err != nil {
		logrus.Fatalf(err.Error())
	}

	go func() {
		for {
			var filename string
			fmt.Scanln(&filename)
			file := fmt.Sprintf("json%s", filename)
			jsonStr, err := ioutil.ReadFile(file)

			if err != nil {
				logrus.Errorf("Error while reading json: %s", err.Error())
			}

			if err := nats.Publish(channelName, jsonStr); err != nil {
				logrus.Errorf("Error while publised msg : %s", err.Error())
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Println("Server shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Error("error while shutting down server", err.Error())
	}

	if err := nats.Close(); err != nil {
		logrus.Error("error while closing channel", err.Error())
	}

}
