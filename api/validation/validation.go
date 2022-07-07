package validation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Field string
	Msg   string
}

func GetValidationMessageForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	}
	return ""
}

func ProcessAndSendValidationErrorMessage(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), GetValidationMessageForTag(fe.Tag())}
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": out})
		return
	}

	c.JSON(http.StatusInternalServerError, "Error during parsing of HTTP request body. Please check it format correctness.")
}