package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	accountSid := os.Getenv("A_SID")
	authToken := os.Getenv("AUTH_TOKEN")
	from := os.Getenv("TWILIO_NUMBER")
	to := "+916379430684"

	message := "Hello, from Mohit today! Twilio is running successfully"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	//Validation(to, "Aayush", client)
	SMS(from, to, message, client)

}

func SMS(from, to, message string, client *twilio.RestClient) {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}

}

func Validation(phone, name string, client *twilio.RestClient) {
	params := &twilioApi.CreateValidationRequestParams{}
	params.SetFriendlyName(name)
	params.SetPhoneNumber(phone)

	resp, err := client.Api.CreateValidationRequest(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.FriendlyName != nil {
			fmt.Println(*resp.FriendlyName)
		} else {
			fmt.Println(resp.FriendlyName)
		}
	}

}
