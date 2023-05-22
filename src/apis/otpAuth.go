package apis

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type PhoneNumber struct {
	PhoneNumber string `json:"phone_number"`
}

type OTP struct {
	OTP string `json:"otp"`
}

func (s Server) SendOTP(c *gin.Context) {
	var phoneNum PhoneNumber

	if err := c.ShouldBind(&phoneNum); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, errBadRequest.Error())
		return
	}

	if err := validatePhoneNumber(phoneNum.PhoneNumber); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	params := &verify.CreateVerificationParams{}
	params = fillCreateVerificationParams(params, phoneNum.PhoneNumber)

	resp, err := s.Client.VerifyV2.CreateVerification(os.Getenv("VERIDY_S_ID"), params)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if resp.Status != nil {
		s.Logger.Infoln(*resp.Status)
		c.SetCookie("PhoneNumber", phoneNum.PhoneNumber, 1200, "/", "localhost", false, true)
		c.JSON(http.StatusOK, "OTP sent successfully")
		return
	}

	s.Logger.Errorln(resp.Status)
	c.JSON(http.StatusInternalServerError, errInternalServer)
}

func (s Server) OTPVerification(c *gin.Context) {
	var otp OTP
	if err := c.ShouldBind(&otp); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, errBadRequest.Error())
		return
	}

	if err := validateOTP(otp.OTP); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	phone, err := returnPhoneFromCookie(c)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	params := &verify.CreateVerificationCheckParams{}
	params = fillCreateVerificationCheckParams(params, phone, otp)

	resp, err := s.Client.VerifyV2.CreateVerificationCheck(os.Getenv("VERIDY_S_ID"), params)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if *resp.Status != "approved" {
		s.Logger.Info(*resp.Status)
		c.JSON(http.StatusBadRequest, *resp.Status)
		return
	}

	c.JSON(http.StatusOK, "logged in successfully")
}

func returnPhoneFromCookie(c *gin.Context) (string, error) {
	phone, err := c.Cookie("PhoneNumber")
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return phone, nil
}

func fillCreateVerificationParams(params *verify.CreateVerificationParams, phone string) *verify.CreateVerificationParams {
	params.SetTo(phone)
	params.SetChannel("sms")

	return params
}

func fillCreateVerificationCheckParams(params *verify.CreateVerificationCheckParams, phone string, otp OTP) *verify.CreateVerificationCheckParams {
	params.SetTo(phone)
	params.SetCode(otp.OTP)

	return params
}

//  func (s Server) enterOTP(c *gin.Context) {
// 	phoneNumber := c.Query("phone_number")
// 	c.HTML(http.StatusOK, "enter_otp.html"gin.H{"phone_number": phoneNumber})
// }
