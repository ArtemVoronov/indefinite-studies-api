//go:build integration
// +build integration

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

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

func CreateUserPutOrPostBody(login any, email any, password any, role any, state any) (string, error) {
	loginField, err := ParseForJsonBody("Login", login)
	if err != nil {
		return "", err
	}
	emailField, err := ParseForJsonBody("Email", email)
	if err != nil {
		return "", err
	}
	passwordField, err := ParseForJsonBody("Password", password)
	if err != nil {
		return "", err
	}
	roleField, err := ParseForJsonBody("Role", role)
	if err != nil {
		return "", err
	}
	stateField, err := ParseForJsonBody("State", state)
	if err != nil {
		return "", err
	}

	userCreateDTO := "{"
	if loginField != "" && emailField != "" && passwordField != "" && roleField != "" && stateField != "" {
		userCreateDTO += loginField + ", " + emailField + ", " + passwordField + ", " + roleField + ", " + stateField
	} else if loginField != "" {
		userCreateDTO += loginField
	} else if emailField != "" {
		userCreateDTO += emailField
	} else if passwordField != "" {
		userCreateDTO += passwordField
	} else if roleField != "" {
		userCreateDTO += roleField
	} else if stateField != "" {
		userCreateDTO += stateField
	}
	userCreateDTO += "}"

	return userCreateDTO, nil
}

func CreateUser(login any, email any, password any, role any, state any) (int, string, error) {
	userCreateDTO, err := CreateUserPutOrPostBody(login, email, password, role, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(userCreateDTO)))
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func GetUser(id string) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/"+id, nil)
	Router.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func GetUsers(limit any, offset any) (int, string, error) {
	queryParams, err := CreateLimitAndOffsetQueryParams(limit, offset)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users"+queryParams, nil)
	Router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func UpdateUser(id any, login any, email any, password any, role any, state any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}
	userUpdateDTO, err := CreateUserPutOrPostBody(login, email, password, role, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/users"+idParam, bytes.NewBuffer([]byte(userUpdateDTO)))
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func DeleteUser(id any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users"+idParam, nil)
	Router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}
