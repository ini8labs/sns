package apis

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/twilio/twilio-go"

	docs "github.com/ini8labs/snsV2/docs"
)

var (
	errBadRequest      error = errors.New("bad request")
	errInvalidPhoneNum error = errors.New("invalid phone number")
	errInvalidOTP      error = errors.New("invalid OTP")
	errIncorrectOTP    error = errors.New("incorrect OTP")
	errInternalServer  error = errors.New("Somethng wrong with the server")
)

type Server struct {
	Logger *logrus.Logger
	Client *twilio.RestClient
}

func NewServer(server Server) error {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.POST("/api/v1/login/otp", server.SendOTP)
	r.POST("/api/v1/login/verify", server.OTPVerification)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r.Run(":8080")
}
