package apis

import (
	"encoding/json"
	"fmt"

	//"github.com/joho/godotenv"

	//"github.com/gin-gonic/gin"
	//twilio "github.com/sfreiberg/gotwilio"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	//"github.com/twilio/twilio-go"
)

func (s Server) SMS(from, to, message string) {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(message)

	resp, err := s.Client.Api.CreateMessage(params)
	if err != nil {
		//fmt.Println("Error sending SMS message: " + err.Error())
		s.Logger.Errorln("Error sending SMS message: " + err.Error())

	} else {
		response, _ := json.Marshal(*resp)
		//fmt.Println("Response: " + string(response))
		s.Logger.Infoln("Response: " + string(response))
	}
}

// Validation() would be depricated
func (s Server) Validation(phone, name string) {
	params := &twilioApi.CreateValidationRequestParams{}
	params.SetFriendlyName(name)
	params.SetPhoneNumber(phone)

	resp, err := s.Client.Api.CreateValidationRequest(params)
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
