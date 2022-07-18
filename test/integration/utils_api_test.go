//go:build integration
// +build integration

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var testHttpClient TestHttpClient = TestHttpClient{}

type TagsApi interface {
	CreateTag(name any, state any) (int, string, error)
	GetTag(id string) (int, string)
	GetTags(limit any, offset any) (int, string, error)
	UpdateTag(id any, name any, state any) (int, string, error)
	DeleteTag(id any) (int, string, error)
}

type TaskApi interface {
	CreateTask(name any, state any) (int, string, error)
	GetTask(id string) (int, string)
	GetTasks(limit any, offset any) (int, string, error)
	UpdateTask(id any, name any, state any) (int, string, error)
	DeleteTask(id any) (int, string, error)
}

type UsersApi interface {
	CreateUser(login any, email any, password any, role any, state any) (int, string, error)
	GetUser(id string) (int, string)
	GetUsers(limit any, offset any) (int, string, error)
	UpdateUser(id any, login any, email any, password any, role any, state any) (int, string, error)
	DeleteUser(id any) (int, string, error)
}

type TestApi interface {
	TagsApi
	TaskApi
	UsersApi
}

type TestHttpClient struct {
}

func (p *TestHttpClient) CreateTask(name any, state any) (int, string, error) {
	taskCreateDTO, err := CreateTaskPutOrPostBody(name, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte(taskCreateDTO)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) GetTask(id string) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+id, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func (p *TestHttpClient) GetTasks(limit any, offset any) (int, string, error) {
	queryParams, err := CreateLimitAndOffsetQueryParams(limit, offset)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks"+queryParams, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) UpdateTask(id any, name any, state any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}
	taskUpdateDTO, err := CreateTaskPutOrPostBody(name, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tasks"+idParam, bytes.NewBuffer([]byte(taskUpdateDTO)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) DeleteTask(id any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/tasks"+idParam, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) CreateTag(name any, state any) (int, string, error) {
	tagCreateDTO, err := CreateTagPutOrPostBody(name, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tags", bytes.NewBuffer([]byte(tagCreateDTO)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) GetTag(id string) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tags/"+id, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func (p *TestHttpClient) GetTags(limit any, offset any) (int, string, error) {
	queryParams, err := CreateLimitAndOffsetQueryParams(limit, offset)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tags"+queryParams, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) UpdateTag(id any, name any, state any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}
	tagUpdateDTO, err := CreateTagPutOrPostBody(name, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tags"+idParam, bytes.NewBuffer([]byte(tagUpdateDTO)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) DeleteTag(id any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/tags"+idParam, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) CreateUser(login any, email any, password any, role any, state any) (int, string, error) {
	userCreateDTO, err := CreateUserPutOrPostBody(login, email, password, role, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(userCreateDTO)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) GetUser(id string) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/"+id, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func (p *TestHttpClient) GetUsers(limit any, offset any) (int, string, error) {
	queryParams, err := CreateLimitAndOffsetQueryParams(limit, offset)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users"+queryParams, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) UpdateUser(id any, login any, email any, password any, role any, state any) (int, string, error) {
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
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) DeleteUser(id any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users"+idParam, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func ParseForJsonBody(paramName string, paramValue any) (string, error) {
	result := ""
	switch paramType := paramValue.(type) {
	case int:
		result = "\"" + paramName + "\": " + strconv.Itoa(paramValue.(int))
	case string:
		result = "\"" + paramName + "\": \"" + paramValue.(string) + "\""
	case nil:
		result = ""
	default:
		return "", fmt.Errorf("unkown type for '%s': %v", paramName, paramType)
	}
	return result, nil
}

func ParseForPathParam(paramName string, paramValue any) (string, error) {
	result := ""
	switch paramType := paramValue.(type) {
	case int:
		result = "/" + strconv.Itoa(paramValue.(int))
	case string:
		result = "/" + paramValue.(string)
	case nil:
		result = ""
	default:
		return "", fmt.Errorf("unkown type for '%s': %v", paramName, paramType)
	}
	return result, nil
}

func ParseForQueryParam(paramName string, paramValue any) (string, error) {
	result := ""
	switch paramType := paramValue.(type) {
	case int:
		result = paramName + "=" + strconv.Itoa(paramValue.(int))
	case nil:
		result = ""
	default:
		return "", fmt.Errorf("unkown type for '%s': %v", paramName, paramType)
	}
	return result, nil
}

func CreateLimitAndOffsetQueryParams(limit any, offset any) (string, error) {
	limitQueryParam, err := ParseForQueryParam("limit", limit)
	if err != nil {
		return "", err
	}
	offsetQueryParam, err := ParseForQueryParam("offset", offset)
	if err != nil {
		return "", err
	}

	queryParams := ""
	if limitQueryParam != "" && offsetQueryParam != "" {
		queryParams += "?" + limitQueryParam + "&" + offsetQueryParam
	} else if limitQueryParam != "" {
		queryParams += "?" + limitQueryParam
	} else if offsetQueryParam != "" {
		queryParams += "?" + offsetQueryParam
	}

	return queryParams, nil
}

func CreateTaskPutOrPostBody(name any, state any) (string, error) {
	nameField, err := ParseForJsonBody("Name", name)
	if err != nil {
		return "", err
	}
	stateField, err := ParseForJsonBody("State", state)
	if err != nil {
		return "", err
	}

	taskCreateDTO := "{"
	if nameField != "" && stateField != "" {
		taskCreateDTO += nameField + ", " + stateField
	} else if nameField != "" {
		taskCreateDTO += nameField
	} else if stateField != "" {
		taskCreateDTO += stateField
	}
	taskCreateDTO += "}"

	return taskCreateDTO, nil
}

func CreateTagPutOrPostBody(name any, state any) (string, error) {
	nameField, err := ParseForJsonBody("Name", name)
	if err != nil {
		return "", err
	}
	stateField, err := ParseForJsonBody("State", state)
	if err != nil {
		return "", err
	}

	tagCreateDTO := "{"
	if nameField != "" && stateField != "" {
		tagCreateDTO += nameField + ", " + stateField
	} else if nameField != "" {
		tagCreateDTO += nameField
	} else if stateField != "" {
		tagCreateDTO += stateField
	}
	tagCreateDTO += "}"

	return tagCreateDTO, nil
}

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
