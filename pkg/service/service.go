package service

import (
	"bytes"
	"encoding/json"
	"flowanalysis/pkg/db"
	"flowanalysis/pkg/inmemory"
	"flowanalysis/pkg/log"
	"flowanalysis/pkg/utils"
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
		Code:    0,
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

	insertID(id)

	endpoint := c.QueryParam("endpoint")
	if endpoint != "" {

		uniquecount := getCount()
		timestamp := time.Now().Format(time.RFC3339)

		statusCode, err := sendUniqueCountAsPost(endpoint, timestamp, uniquecount)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed: could not send data to endpoint, error: %v", err))
		}

		log.Print(log.Info, "HTTP POST request to endpoint [%s] responded with status code: %d", endpoint, statusCode)
	}

	return c.String(http.StatusOK, "ok")
}

func sendUniqueCountAsPost(endpoint string, timestamp string, count int) (int, error) {
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	payload := map[string]interface{}{
		"timestamp":    timestamp,
		"unique_count": count,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, &ApiError{
			ErrType: "Request Error",
			Message: "Failed to marshal JSON payload",
			Code:    0,
		}
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, &ApiError{
			ErrType: "Request Error",
			Message: err.Error(),
			Code:    0,
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
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

func insertID(id string) {
	if utils.UseInMemory() {
		activeMap := inmemory.GetActiveFlowMap()
		activeMap.Set(id, "1")
	} else {
		InsertEntry(id, time.Now())
	}

}

func getCount() int {
	if utils.UseInMemory() {
		activeMap := inmemory.GetActiveFlowMap()
		return activeMap.Count()
	}
	count, _ := CountEntries(time.Now())
	return count
}

func InsertEntry(id string, timestamp time.Time) {
	instance := db.GetPostgreDbInstance()
	err := instance.Connect()
	if err != nil {
		log.Print(log.Warn, "Failed to connect to the database: %v", err)
	}
	defer instance.Disconnect()

	formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")

	query := `
		INSERT INTO my_table (unique_id, timestamp)
		VALUES ($1, $2)
		ON CONFLICT (unique_id)
		DO UPDATE SET timestamp = EXCLUDED.timestamp;`
	instance.Insert(query, id, formattedTimestamp)

}

func CountEntries(inputTime time.Time) (int, error) {
	instance := db.GetPostgreDbInstance()
	err := instance.Connect()
	if err != nil {
		log.Print(log.Warn, "Failed to connect to the database: %v", err)
	}

	startTime := inputTime.Add(-1 * time.Minute)

	endTime := inputTime

	startTimeStr := startTime.Format("2006-01-02 15:04:05")
	endTimeStr := endTime.Format("2006-01-02 15:04:05")
	log.Print(log.Info, "Start time %s and end time for count unique entries query %s", startTime, endTime)
	// Construct the SQL query
	query := `
		SELECT COUNT(*)
		FROM my_table
		WHERE timestamp >= $1 AND timestamp < $2;`

	result, err := instance.Query(query, startTimeStr, endTimeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}

	rows, ok := result.([]map[string]interface{})
	if !ok || len(rows) == 0 {
		return 0, fmt.Errorf("unexpected result format or no rows returned")
	}
	count, ok := rows[0]["count"].(int64)
	if !ok {
		return 0, fmt.Errorf("failed to parse count from query result")
	}

	return int(count), nil
}
