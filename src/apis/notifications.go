package apis

import (
	"encoding/json"
	"os"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func (s Server) SMS(to, message string) (string, error) {
	if err := validatePhoneNumber(to); err != nil {
		s.Logger.Error(err)
		return "", err
	}

	params := &twilioApi.CreateMessageParams{}
	params = fillCreateMessageParams(params, to, message)

	resp, err := s.Client.Api.CreateMessage(params)
	if err != nil {
		s.Logger.Errorln("Error sending SMS message: " + err.Error())
		return "", err
	}

	response, _ := json.Marshal(*resp)
	s.Logger.Infoln("Response: " + string(response))
	return string(response), nil
}

func fillCreateMessageParams(params *twilioApi.CreateMessageParams, to, message string) *twilioApi.CreateMessageParams {
	params.SetTo(to)
	params.SetFrom(os.Getenv("TWILIO_NUMBER"))
	params.SetBody(message)

	return params
}
