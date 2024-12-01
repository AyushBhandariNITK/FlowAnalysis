package pkg

import (
	"flowanalysis/pkg/log"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type ApiError struct {
	Code    int
	Message string
	ErrType string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("[%s] %d: %s", e.ErrType, e.Code, e.Message)
}

var (
	ErrConnectionRefused = &ApiError{
		ErrType: "Connection Refused",
		Message: "The server refused the connection.",
		Code:    0, // No HTTP status code
	}
	ErrInvalidEndpoint = &ApiError{
		ErrType: "Invalid Endpoint",
		Message: "The API endpoint was not found (404).",
		Code:    404,
	}
	ErrBadRequest = &ApiError{
		ErrType: "Bad Request",
		Message: "The request was malformed (400).",
		Code:    400,
	}
)

func GetAcceptHandler(c echo.Context) error {
	id := c.QueryParam("id")

	if id == "" {
		return c.String(http.StatusBadRequest, "failed")
	}
	activeMap.Set(id, "1")
	endpoint := c.QueryParam("endpoint")
	if endpoint != "" {
		uniqueCount := activeMap.Count()
		statusCode, err := sendUniqueCountAsQuery(endpoint, uniqueCount)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed: could not send data to endpoint, error: %v", err))
		}

		log.Print(log.Info, "HTTP GET request to endpoint [%s] responded with status code: %d\n", endpoint, statusCode)

	}
	return c.String(http.StatusOK, "ok")
}

func sendUniqueCountAsQuery(endpoint string, uniqueCount int) (int, error) {
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	url := fmt.Sprintf("%s?unique_count=%d", endpoint, uniqueCount)
	resp, err := client.Get(url)
	if err != nil {
		if isConnectionRefusedError(err) {
			return 0, ErrConnectionRefused
		}
		return 0, &ApiError{
			ErrType: "Request Error",
			Message: err.Error(),
			Code:    0,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return resp.StatusCode, ErrInvalidEndpoint
	} else if resp.StatusCode == 400 {
		return resp.StatusCode, ErrBadRequest
	}

	return resp.StatusCode, nil
}

func isConnectionRefusedError(err error) bool {
	return strings.Contains(err.Error(), "connection refused")
}
