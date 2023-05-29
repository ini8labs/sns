package apis

import (
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/twilio/twilio-go"

	mw "github.com/ini8labs/sns/src/middleware"
)

var (
	twilioNumber       string
	serviceID          string
	errBadRequest      error = errors.New("bad request")
	errInvalidPhoneNum error = errors.New("invalid phone number")
	errInvalidOTP      error = errors.New("invalid OTP")
	errIncorrectOTP    error = errors.New("incorrect OTP")
	errInternalServer  error = errors.New("Somethng wrong with the server")
	errInvalidToken    error = errors.New("invalid token")
	errUnauthorized    error = errors.New("unauthorized")
	errPayload         error = errors.New("payload not found")
)

const (
	twilioNumberEnv = "TWILIO_NUMBER"
	serviceIDEnv    = "VERIDY_S_ID"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("couldn't load env")
		panic(err.Error())
	}

	twilioNumber = os.Getenv(twilioNumberEnv)
	serviceID = os.Getenv(serviceIDEnv)

	if twilioNumber == "" || serviceID == "" {
		log.Fatal("Twilio Number and ServiceID are not set")
		panic("Twilio environment variables are not set")
	}
}

type Server struct {
	Logger *logrus.Logger
	Client *twilio.RestClient
}

func NewServer(server Server) error {
	r := gin.Default()
	//docs.SwaggerInfo.BasePath = "/api/v1"

	r.POST("/api/v1/login/otp", server.SendOTP)
	r.POST("/api/v1/login/verify", server.OTPVerification)
	r.POST("api/v1/notify", mw.AuthMiddleware([]byte(os.Getenv("A_SID"))), server.SMSCheck)
	r.POST("api/v1/token/refresh", mw.AuthMiddleware([]byte(os.Getenv("A_SID"))), mw.JwtrefreshToken)

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r.Run(":8080")
}
