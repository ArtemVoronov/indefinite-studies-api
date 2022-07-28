package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/validation"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/app/utils"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// TODO: finish implementation

var hmacSecret []byte
var accessTokenDuration time.Duration
var refreshTokenDuration time.Duration
var tokenIssuer string
var once sync.Once

func Setup() {
	once.Do(func() {
		hmacSecret = createHmacSecret()
		accessTokenDuration = createAccessTokenDuration()
		refreshTokenDuration = createRefreshTokenDuration()
		tokenIssuer = createIssuer()
	})
}

type CredentialsValidationResult struct {
	userId  int
	isValid bool
}

type AuthenicationDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerificationDTO struct {
	Token string `json:"token" binding:"required"`
}

type AuthenicationResultDTO struct {
	AccessToken           string `json:"accessToken" binding:"required"`
	RefreshToken          string `json:"refreshToken" binding:"required"`
	AccessTokenExpiredAt  string `json:"accessTokenExpiredAt" binding:"required"`
	RefreshTokenExpiredAt string `json:"refreshTokenExpiredAt" binding:"required"`
}

type UserClaims struct {
	Email string
	jwt.RegisteredClaims
}

func Authenicate(c *gin.Context) {
	var authenicationDTO AuthenicationDTO

	if err := c.ShouldBindJSON(&authenicationDTO); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	// TODO: add counter of invalid athorizations, then use it for temporary blocking access
	validatoionResult, err := checkUserCredentials(authenicationDTO.Email, authenicationDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error")
		log.Printf("error during authenication: %v\n", err)
		return
	}
	if !validatoionResult.isValid || validatoionResult.userId == -1 {
		c.JSON(http.StatusBadRequest, api.ERROR_WRONG_PASSWORD_OR_EMAIL)
		return
	}

	expireAtForAccessToken := jwt.NewNumericDate(time.Now().Add(accessTokenDuration))
	expireAtForRefreshToken := jwt.NewNumericDate(time.Now().Add(refreshTokenDuration))

	accessToken, err := createToken(expireAtForAccessToken, authenicationDTO.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to authenicate")
		log.Printf("error during authenication: %v\n", err)
		return
	}

	refreshToken, err := createToken(expireAtForRefreshToken, authenicationDTO.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to authenicate")
		log.Printf("error during authenication: %v\n", err)
		return
	}

	result := &AuthenicationResultDTO{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  expireAtForAccessToken.String(),
		RefreshTokenExpiredAt: expireAtForRefreshToken.String(),
	}

	err = db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
		err := queries.CreateRefreshToken(tx, ctx, validatoionResult.userId, refreshToken, expireAtForRefreshToken.Time)
		return err
	})()

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error")
		log.Printf("error during authenication: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func Verify(c *gin.Context) {
	var verificationDTO VerificationDTO

	if err := c.ShouldBindJSON(&verificationDTO); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	token, err := jwt.ParseWithClaims(verificationDTO.Token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		fmt.Println("---------------------------------------------")
		fmt.Printf("%v %v\n", claims.Email, claims.RegisteredClaims)
		fmt.Println("---------------------------------------------")
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

func checkUserCredentials(email string, password string) (CredentialsValidationResult, error) {
	var result CredentialsValidationResult

	data, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
		userId, isValid, err := queries.IsValidCredentials(tx, ctx, email, utils.CreateSHA512HashHexEncoded(password))
		return CredentialsValidationResult{userId: userId, isValid: isValid}, err
	})()

	if err != nil && err != sql.ErrNoRows {
		return result, fmt.Errorf("unable to check credentials : %s", err)
	}

	result, ok := data.(CredentialsValidationResult)
	if !ok {
		return result, fmt.Errorf("unable to check credentials : %s", api.ERROR_ASSERT_RESULT_TYPE)
	}

	return result, nil
}

func createToken(expireAt *jwt.NumericDate, email string) (string, error) {
	claims := UserClaims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: expireAt,
			Issuer:    tokenIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString(hmacSecret)
	return signedToken, err
}

func createHmacSecret() []byte {
	sign, signExists := os.LookupEnv("JWT_SIGN")
	if !signExists {
		log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", "JWT_SIGN")
	}
	return []byte(sign)
}

func createAccessTokenDuration() time.Duration {
	duration, durationExists := os.LookupEnv("JWT_ACCESS_DURATION_IN_MINUTES")
	if !durationExists {
		log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", "JWT_ACCESS_DURATION_IN_MINUTES")
	}

	minutes, err := strconv.Atoi(duration)
	if err != nil {
		log.Fatalf("Wrong value of environment variable: %s. It should be integer number", "JWT_ACCESS_DURATION_IN_MINUTES")
	}
	return time.Minute * time.Duration(minutes)
}

func createRefreshTokenDuration() time.Duration {
	duration, durationExists := os.LookupEnv("JWT_REFRESH_DURATION_IN_HOURS")
	if !durationExists {
		log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", "JWT_REFRESH_DURATION_IN_HOURS")
	}

	hours, err := strconv.Atoi(duration)
	if err != nil {
		log.Fatalf("Wrong value of environment variable: %s. It should be integer number", "JWT_REFRESH_DURATION_IN_HOURS")
	}
	return time.Hour * time.Duration(hours)
}

func createIssuer() string {
	issuer, issuerExists := os.LookupEnv("JWT_ISSUER")
	if !issuerExists {
		log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", "JWT_ISSUER")
	}
	return issuer
}
