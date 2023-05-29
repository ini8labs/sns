package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mw "github.com/ini8labs/sns/src/middleware"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type PhoneNumber struct {
	PhoneNumber string `json:"phone_number"`
}

type OTP struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
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

	params := fillCreateVerificationParams(phoneNum)
	resp, err := s.Client.VerifyV2.CreateVerification(serviceID, params)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if resp.Status != nil {
		s.Logger.Infoln(*resp.Status)
		c.JSON(http.StatusOK, "OTP sent successfully")
		return
	}

	s.Logger.Errorln(resp.Status)
	c.JSON(http.StatusInternalServerError, errInternalServer)
}

func (s Server) OTPVerification(c *gin.Context) {
	var otp OTP
	if err := c.ShouldBind(&otp); err != nil {
		c.JSON(http.StatusBadRequest, errBadRequest)
		s.Logger.Error(err)

		return
	}

	if err := validateOTP(otp.OTP); err != nil {
		c.JSON(http.StatusBadRequest, errBadRequest)
		s.Logger.Error(err)

		return
	}

	params := fillCreateVerificationCheckParams(otp)
	resp, err := s.Client.VerifyV2.CreateVerificationCheck(serviceID, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, errBadRequest)
		s.Logger.Error(err)

		return
	}

	if *resp.Status != "approved" {
		c.JSON(http.StatusBadRequest, *resp.Status)
		s.Logger.Info(*resp.Status)

		return
	}

	if _, err := mw.JwtCreateToken(c, otp.PhoneNumber); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, errInternalServer.Error())

		return
	}

	c.JSON(http.StatusOK, "logged in successfully")
}

func fillCreateVerificationParams(phone PhoneNumber) *verify.CreateVerificationParams {
	params := &verify.CreateVerificationParams{}
	params.SetTo(phone.PhoneNumber)
	params.SetChannel("sms")

	return params
}

func fillCreateVerificationCheckParams(OTPAuth OTP) *verify.CreateVerificationCheckParams {
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(OTPAuth.PhoneNumber)
	params.SetCode(OTPAuth.OTP)

	return params
}

//  func (s Server) enterOTP(c *gin.Context) {
// 	phoneNumber := c.Query("phone_number")
// 	c.HTML(http.StatusOK, "enter_otp.html"gin.H{"phone_number": phoneNumber})
// }
