package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func decodePayload(r *http.Request) (string, []any, error) {
	var payload map[string]any
	err := json.NewDecoder(r.Body).Decode(&payload)
	defer r.Body.Close()
	if err != nil {
		return "", nil, fmt.Errorf("invalid json format")
	}
	query, ok := payload["query"].(string)
	if !ok {
		return "", nil, fmt.Errorf("invalid json body: no query")
	}
	args, ok := payload["args"].([]any)
	if !ok {
		args = nil
	}
	return query, args, nil
}

type queryHandler struct {
	env *env
}

func QueryHandler(env *env) *queryHandler {
	return &queryHandler{env: env}
}

func (t *queryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query, args, err := decodePayload(r)
	if err != nil {
		ErrorHandler("bad request", err.Error(), http.StatusBadRequest).ServeHTTP(w, r)
		return
	}
	result, err := t.env.database.Query(query, args...)
	if err != nil {
		ErrorHandler("bad request", "failed to query data: "+err.Error(), http.StatusBadRequest).ServeHTTP(w, r)
		return
	}
	responseBody, err := json.Marshal(result)
	if err != nil {
		ErrorHandler("internal server error", "failed encoding response body", http.StatusInternalServerError).ServeHTTP(w, r)
		return
	}
	w.Write(responseBody)
}

type execHandler struct {
	env *env
}

func ExecHandler(env *env) *execHandler {
	return &execHandler{env: env}
}

func (t *execHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query, args, err := decodePayload(r)
	r.Body.Close()
	if err != nil {
		ErrorHandler("bad request", "invalid json body", http.StatusBadRequest).ServeHTTP(w, r)
		return
	}
	err = t.env.database.Exec(query, args)
	if err != nil {
		ErrorHandler("bad request", "failed to query data", http.StatusBadRequest).ServeHTTP(w, r)
		return
	}
	w.Write([]byte{})
}

type apiError struct {
	Error   string
	Message string
}

type errorHandler struct {
	Error   string
	Message string
	Code    int
}

func ErrorHandler(err string, message string, code int) *errorHandler {
	return &errorHandler{Error: err, Message: message, Code: code}
}

func (t *errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(t.Code)
	json.NewEncoder(w).Encode(apiError{Error: t.Error, Message: t.Message})

}
