//go:build integration
// +build integration

package integration

import (
	"fmt"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
)

var (
	ERROR_USER_LOGIN_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Login\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_USER_EMAIL_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Email\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_USER_PASSWORD_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Password\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_USER_ROLE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Role\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_USER_STATE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_USER_ALL_ARE_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Login\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"Email\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"Password\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"Role\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_USER_CREATE_STATE_WRONG_VALUE string = fmt.Sprintf("Unable to create user. Wrong 'State' value. Possible values: %v", entities.GetPossibleUserStates())
	ERROR_USER_UPDATE_STATE_WRONG_VALUE string = fmt.Sprintf("Unable to update user. Wrong 'State' value. Possible values: %v", entities.GetPossibleUserStates())
)
