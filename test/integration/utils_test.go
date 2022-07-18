//go:build integration
// +build integration

package integration

import (
	"fmt"
	"strconv"
)

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
