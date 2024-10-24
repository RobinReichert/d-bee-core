package dbee

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func newBodyReader(query string, args ...any) (io.Reader, error) {
	requestBody := map[string]any{
		"query": query,
		"args":  args,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	return bytes.NewReader(jsonData), nil
}

type Connection interface {
	Query(query string, args ...any) ([]map[string]any, error)
	Exec(query string, args ...any) error
}

type connection struct {
	baseUrl string
}

func Connect(baseUrl string) *connection {
	log.Println("connecting to d-bee api")
	return &connection{baseUrl: baseUrl}
}

func (t *connection) Query(query string, args ...any) ([]map[string]any, error) {
	bodyReader, err := newBodyReader(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create body reader: %w", err)
	}
	response, err := http.Post(t.baseUrl+"/query", "application/json", bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed post request: %w", err)
	}
	defer response.Body.Close()
	if response.Header.Get("Content-Type") != "application/json" {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body with wrong content type: %w", err)
		}
		return nil, errors.New("unexpected response format, message: " + string(body))
	}
	if response.StatusCode == http.StatusOK {
		var responseBody []map[string]any
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			return nil, fmt.Errorf("failed to decode response body: %w", err)
		}
		return responseBody, nil
	} else {
		var errorBody map[string]any
		err = json.NewDecoder(response.Body).Decode(&errorBody)
		if err != nil {
			return nil, fmt.Errorf("failed to decode error: %w", err)
		}
		if msg, ok := errorBody["message"].(string); ok {
			return nil, errors.New(msg)
		}
		return nil, errors.New("failed to find message in error")
	}
}

func (t *connection) Exec(query string, args ...any) error {
	bodyReader, err := newBodyReader(query, args...)
	if err != nil {
		return fmt.Errorf("failed to create body reader: %w", err)
	}
	response, err := http.Post(t.baseUrl+"/exec", "application/json", bodyReader)
	if err != nil {
		return fmt.Errorf("failed to post query: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		var errorBody map[string]any
		err = json.NewDecoder(response.Body).Decode(&errorBody)
		if err != nil {
			return fmt.Errorf("failed to decode error: %w", err)
		}
		if msg, ok := errorBody["message"].(string); ok {
			return errors.New(msg)
		}
		return errors.New("error wrong format")
	}
	return nil
}
