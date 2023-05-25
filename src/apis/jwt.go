package apis

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type PayLoad struct {
	userName string
	userID   string
	phone    string
}

var (
	signingKey = []byte(os.Getenv("A_SID"))
	jwtToken   = "jwt_token"
)

func (s Server) JwtSetToken(c *gin.Context, userID, username, phone string) error {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   userID,
		"userName": username,                         //TODO: get username from lsdb.getUserInfoByGovid
		"phone":    phone,                            //TODO:  get phone form lsdb.getUserInfoByGovid
		"exp":      time.Now().Add(time.Hour).Unix(), //TODO: check if it works with UTC time
	})

	// delete "PhoneNumber" cookie
	if err := s.unsetCookie(c, "PhoneNumber"); err != nil {
		c.JSON(http.StatusInternalServerError, errInternalServer)

		return err
	}

	// Sign the token with the secret key
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		s.Logger.Error("Failed to generate token")
		c.JSON(http.StatusInternalServerError, errInternalServer)

		return err
	}

	// Set the token as a cookie on the client side
	c.SetCookie(jwtToken, tokenString, int(time.Hour.Seconds()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, "Token generated successfully")
	return nil
}

func (s Server) JwtTokenAuth(c *gin.Context) (PayLoad, error) {
	// Retrieve the token from the cookie
	tokenString, err := c.Cookie(jwtToken)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusUnauthorized, errUnauthorized.Error())

		return PayLoad{}, nil
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return signingKey, nil
	})

	if err != nil {
		s.Logger.Error(errInvalidToken)
		c.JSON(http.StatusUnauthorized, errInvalidToken.Error())
		c.Abort()
		return PayLoad{}, err
	}

	payLoad, err := jwtPayload(token)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusUnauthorized, errUnauthorized.Error())
		return PayLoad{}, err
	}

	//c.JSON(http.StatusOK, payLoad)
	return payLoad, nil
}

func jwtPayload(token *jwt.Token) (PayLoad, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var payLoad PayLoad
		payLoad.userName = claims["userName"].(string)
		payLoad.userID = claims["userID"].(string)
		payLoad.phone = claims["phone"].(string)

		return payLoad, nil
	}
	return PayLoad{}, errPayload
}

func (s Server) unsetCookie(c *gin.Context, cookieName string) error {
	c.SetCookie(cookieName, "", -1, "/", "", false, true)

	// check to see if cookie still exists
	if err := c.ShouldBindHeader(cookieName); err != nil {
		s.Logger.Errorf("Failed to expire cookie: %v", err)

		return err
	}

	return nil
}
