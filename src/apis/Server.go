package apis

import (
	"errors"
	"sns/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/twilio/twilio-go"
)

var (
	errInvalidPhoneNum error = errors.New("invalid phone number")
	errInvalidOTP      error = errors.New("invalid OTP")
)

type Server struct {
	*logrus.Logger
	Client *twilio.RestClient
}

func NewServer(server Server) error {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.POST("/api/v1/Login/OTP", server.sendOTP)
	r.GET("/enter-otp", server.enterOTP)
	r.POST("/api/v1/Login/verify", server.OTPVerification)
	r.GET("api/v1/push", server.pushNotification)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r.Run(":8080")
}
