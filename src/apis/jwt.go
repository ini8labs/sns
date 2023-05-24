package apis

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	signingKey = []byte("your-secret-key") //TODO: Change this to your own secret key or API key
	jwtToken   = "jwt_token"
)

func (s Server) SetToken(c *gin.Context, username, phone string) error {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": username,                         //TODO: get username from lsdb.getUserInfoByGovid
		"phone":  phone,                            //TODO:  get phone form lsdb.getUserInfoByGovid
		"exp":    time.Now().Add(time.Hour).Unix(), //TODO: check if it works with UTC time
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

func (s Server) authMiddleware(c *gin.Context) {
	// Retrieve the token from the cookie
	tokenString, err := c.Cookie(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return signingKey, nil
	})

	if err != nil {
		s.Logger.Error(errInvalidToken)
		c.JSON(http.StatusUnauthorized, errInvalidToken.Error())
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract session information from the claims
		username := claims["username"].(string)

		// Set the session information in the Gin context
		c.Set("username", username)
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
	}
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
