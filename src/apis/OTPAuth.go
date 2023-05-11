package apis

import (
	"net/http"
	"os"

	//"github.com/sirupsen/logrus"
	//twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/gin-gonic/gin"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type PhoneNumber struct {
	PhoneNumber string `json:"phone_number"`
}

type OTP struct {
	OTP string `json:"otp"`
}

type Verification struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
}

func (s Server) sendOTP(c *gin.Context) {
	var phoneNum PhoneNumber

	if err := c.ShouldBind(&phoneNum); err != nil {
		s.Logger.Error(errInvalidPhoneNum)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	if err := validatePhoneNumber(phoneNum.PhoneNumber); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	params := &verify.CreateVerificationParams{}
	params.SetTo(phoneNum.PhoneNumber)
	params.SetChannel("sms")

	resp, err := s.Client.VerifyV2.CreateVerification(os.Getenv("VERIDY_S_ID"), params)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, "something wrong with server")
		return
	} else {
		if resp.Status != nil {
			s.Logger.Infoln(*resp.Status)
			c.JSON(http.StatusAccepted, "OTP sent successfully")
		} else {
			s.Logger.Infoln(resp.Status)
			c.JSON(http.StatusForbidden, "some 403")
		}
	}
}

func (s Server) OTPVerification(c *gin.Context) {
	var verification Verification
	if err := c.ShouldBind(&verification); err != nil {
		s.Logger.Error(errInvalidOTP)
		c.JSON(http.StatusBadRequest, errInvalidOTP.Error())
		return
	}

	if err := validateOTP(verification.OTP); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(verification.PhoneNumber)
	params.SetCode(verification.OTP)

	resp, err := s.Client.VerifyV2.CreateVerificationCheck("VAbf0905e98d803fef6df430b2197710ee", params)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, errInvalidOTP.Error())
		return
	} else {
		if *resp.Status != "approved" {
			s.Logger.Info(*resp.Status)
			c.JSON(http.StatusUnauthorized, "incorrect OTP")
			return
		}
	}
	c.JSON(http.StatusOK, "logged in")
}

func (s Server) enterOTP(c *gin.Context) {
	phoneNumber := c.Query("phone_number")
	c.HTML(http.StatusOK, "enter_otp.html", gin.H{"phone_number": phoneNumber})
}
