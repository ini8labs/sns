package apis

import (
	"encoding/json"
	"fmt"

	//"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	twilio "github.com/sfreiberg/gotwilio"
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
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}

}

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

// func (s Server) sendVerificationCode(to, from string) error {
//     verificationParams := &client.CreateVerificationParams{
//         To:   to,
//         Channel: "sms",
//         From: from,
//     }

//     _, err := client.Verify.CreateVerification(verificationParams)
//     return err
// }

func (s Server) pushNotification(c *gin.Context) {
	from := "+12707477263"
	to := "+919944105595"
	body := "Hello, this is a push SMS!"

	twilio.sms

	_, exception, err := twilio.SendSMS(from, to, body, "", "")
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	if exception != nil {
		fmt.Println("Exception:", exception.Message)
		return
	}

	fmt.Println("SMS sent successfully!")

}
