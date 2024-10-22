package dbee

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
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
	client  *http.Client
}

func Connect(baseUrl string) *connection {
	return &connection{baseUrl: baseUrl, client: &http.Client{
		Timeout: 10 * time.Second,
	}}
}

func (t *connection) Query(query string, args ...any) ([]map[string]any, error) {
	bodyReader, err := newBodyReader(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create body reader: %w", err)
	}
	request, err := http.NewRequest("POST", t.baseUrl+"/query", bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	response, err := t.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed post request: %w", err)
	}
	defer response.Body.Close()
	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
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
		return nil, fmt.Errorf("internal error")
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
