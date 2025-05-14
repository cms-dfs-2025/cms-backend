package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

var ErrAuthFormat = errors.New("cms: incorrect auth format")

func decodeBase64(encoded []byte) ([]byte, error) {
	var decoded = make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
	n, err := base64.StdEncoding.Decode(decoded, encoded)

	if err != nil {
		return decoded, err
	}

	decoded = decoded[:n]
	return decoded, nil
}

// parses base64(base64(handle):password)
// returns handle, password, error
func parseBasicAuthorization(header_value string) (string, string, error) {
	if len(header_value) < 6 {
		return "", "", ErrAuthFormat
	}

	basic := header_value[:6]

	if basic != "Basic " {
		return "", "", ErrAuthFormat
	}

	passBytes, err := decodeBase64([]byte(header_value[6:]))

	if err != nil {
		return "", "", err
	}

	colon_idx := slices.Index(passBytes, ':')

	if colon_idx == -1 {
		return "", "", ErrAuthFormat
	}

	handleBytes, err := decodeBase64(passBytes[:colon_idx])

	if err != nil {
		return "", "", err
	}

	passwordBytes := passBytes[(colon_idx + 1):]

	return string(handleBytes), string(passwordBytes), nil
}

type signupBody struct {
	Handle   *string `json:"handle" binding:"required"`
	Password *string `json:"password" binding:"required"`
}

func (handler ServerContext) HandleSignup(c *gin.Context) {
	var body signupBody
	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = CreateUser(*body.Handle, *body.Password, false, handler.db)

	if err != nil {
		if err == ErrUserExists {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "User already exists"})
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Internal error"})
		}
		return
	}

	c.Status(http.StatusOK)
	return
}

func (handler ServerContext) LoginMiddleware(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	handle, password, err := parseBasicAuthorization(auth)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"message": "Basic auth error"})
		return
	}

	user, err := AuthorizeUser(handle, password, handler.db)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"message": "Unauthorized"})
		return
	}

	c.Set("user", user)
}

func (handler ServerContext) HandleLogin(c *gin.Context) {
	c.Status(http.StatusOK)
	return
}

type changePwBody struct {
	Password *string `json:"password" binding:"required"`
}

func (handler ServerContext) HandleChangePw(c *gin.Context) {
	user := c.MustGet("user").(UserRow)

	var body changePwBody
	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = ChangeUserPassword(user.Handle, *body.Password, handler.db)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	return
}
