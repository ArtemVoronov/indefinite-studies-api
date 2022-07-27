package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/validation"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// TODO: finish implementation

var hmacSecret []byte
var tokenDuration time.Duration
var tokenIssuer string
var once sync.Once

func Setup() {
	once.Do(func() {
		hmacSecret = createHmacSecret()
		tokenDuration = createTokenDuration()
		tokenIssuer = createIssuer()
	})
}

type AuthenicationDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerificationDTO struct {
	Token string `json:"token" binding:"required"`
}

type AuthenicationResultDTO struct {
	Token     string `json:"token" binding:"required"`
	ExpiredAt string `json:"expiredAt" binding:"required"`
}

type MyCustomClaims struct {
	Email string
	jwt.RegisteredClaims
}

func Authenicate(c *gin.Context) {
	var authenicationDTO AuthenicationDTO

	if err := c.ShouldBindJSON(&authenicationDTO); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	isValid := checkUserCredentials(authenicationDTO.Email, authenicationDTO.Password)
	if !isValid {
		c.JSON(http.StatusBadRequest, "Wrong password or email")
		return
	}

	expireTime := jwt.NewNumericDate(time.Now().Add(tokenDuration))

	claims := MyCustomClaims{
		authenicationDTO.Email,
		jwt.RegisteredClaims{
			ExpiresAt: expireTime,
			Issuer:    tokenIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(hmacSecret)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to authenicate")
		fmt.Printf("error during authenication: %v\n", err)
		return
	}

	result := &AuthenicationResultDTO{Token: signedToken, ExpiredAt: expireTime.String()}
	c.JSON(http.StatusOK, result)
}

func Verify(c *gin.Context) {
	var verificationDTO VerificationDTO

	if err := c.ShouldBindJSON(&verificationDTO); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	token, err := jwt.ParseWithClaims(verificationDTO.Token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired") {
			c.JSON(http.StatusOK, false)
			return
		}
		c.JSON(http.StatusInternalServerError, "Unable to verify")
		return
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		fmt.Printf("%v %v", claims.Email, claims.RegisteredClaims)
	} else {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, token.Valid)
}

func Refresh(c *gin.Context) {
	// TODO:
}

func ForceExpireToken(c *gin.Context) {
	// TODO:
}

func checkUserCredentials(email string, password string) bool {
	// TODO implement
	return true
}

func createHmacSecret() []byte {
	sign, signExists := os.LookupEnv("JWT_SIGN")
	if !signExists {
		log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", "JWT_SIGN")
	}
	return []byte(sign)
}

func createTokenDuration() time.Duration {
	duration, durationExists := os.LookupEnv("JWT_DURATION_IN_MINUTES")
	if !durationExists {
		log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", "JWT_DURATION_IN_MINUTES")
	}

	minutes, err := strconv.Atoi(duration)
	if err != nil {
		log.Fatalf("Wrong value of environment variable: %s. It should be integer number", "JWT_DURATION_IN_MINUTES")
	}
	return time.Minute * time.Duration(minutes)
}

func createIssuer() string {
	issuer, issuerExists := os.LookupEnv("JWT_ISSUER")
	if !issuerExists {
		log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", "JWT_ISSUER")
	}
	return issuer
}
