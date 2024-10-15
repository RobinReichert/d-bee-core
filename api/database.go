package dbee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type connection struct {
	baseUrl string
}

func ConnectToDBEE(baseUrl string) connection {
	return connection{baseUrl: baseUrl}
}

func (t *connection) Query(query string, args ...any) (map[string]any, error) {
	requestBody := map[string]any{
		"query": query,
		"args":  args,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	response, err := http.Post(t.baseUrl+"query", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to post query: %w", err)
	}
	switch response.StatusCode {
	case http.StatusOK:
		var responseBody map[string]any
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			return nil, fmt.Errorf("failed to decode response body: %w", err)
		}
		return responseBody, nil
	case http.StatusBadRequest:
		return nil, fmt.Errorf("bad request")
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("internal server error")
	}
	return nil, fmt.Errorf("unexpected status code")
}
