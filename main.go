package main

import (
	//"encoding/json"
	//"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ini8labs/sns/src/apis"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"
	//verify "github.com/twilio/twilio-go/rest/verify/v2"
)

// @title My API
// @version 1.0
// @description This is Lottery SMS Notification Service API
// @host localhost:3000
// @BasePath /api/v1
// @schemes http
func main() {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	accountSid := os.Getenv("A_SID")
	authToken := os.Getenv("AUTH_TOKEN")

	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	server := apis.Server{
		Logger: logger,
		Client: twilioClient,
	}

	go func() {
		if err := apis.NewServer(server); err != nil {
			panic(err)
		}
	}()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt
	logger.Info("Closing the Server")
}
