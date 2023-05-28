package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	signingKey = []byte(os.Getenv("A_SID"))
	jwtToken   = "jwt_token"
)

func JwtCreateToken(c *gin.Context, phone string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": phone,
		"exp":   time.Now().Add(time.Hour).Unix(), //TODO: check if it works with UTC time
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		logrus.Error("Failed to generate token")

		return "", err
	}

	c.SetCookie(jwtToken, tokenString, int(time.Hour.Seconds()), "/", "localhost", false, true)
	logrus.Info("Token generated successfully")

	return jwtToken, nil
}

func JwtIsAuthorized(c *gin.Context) (bool, error) {
	tokenString, err := c.Cookie(jwtToken)
	if err != nil {
		logrus.Error(err)

		return false, nil
	}

	// Parse and validate the token
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Error(jwt.ErrSignatureInvalid.Error())

			return nil, jwt.ErrSignatureInvalid
		}

		return signingKey, nil
	})

	if err != nil {
		logrus.Error(err)

		return false, err
	}

	logrus.Info("Authorized")

	return true, nil
}

func JwtRefreshAccessToken(refreshTokenString string) (string, error) {

	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Error(jwt.ErrSignatureInvalid.Error())

			return nil, jwt.ErrSignatureInvalid
		}

		return signingKey, nil
	})

	// Extract phone numebr from the refresh token claims
	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	phone, ok := claims["Phone"].(string)
	if !ok {
		return "", err
	}

	// Create a new access token with updated expiration time
	accessTokenClaims := jwt.MapClaims{
		"Phone": phone,
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedAccessToken, err := accessToken.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

func JwtrefreshToken(c *gin.Context) {
	// Get the refresh token from the request
	refreshTokenString, err := c.Cookie(jwtToken)
	if err != nil {

		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, "Unauthorized token")
	}

	// Refresh the access token
	newAccessToken, err := JwtRefreshAccessToken(refreshTokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
