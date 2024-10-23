package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func decodePayload(r *http.Request) (string, []any, error) {
	var payload map[string]any
	err := json.NewDecoder(r.Body).Decode(&payload)

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
	defer r.Body.Close()
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
		ErrorHandler("server error", "failed encoding response body", http.StatusInternalServerError).ServeHTTP(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println(err)
	}

}

type execHandler struct {
	env *env
}

func ExecHandler(env *env) *execHandler {
	return &execHandler{env: env}
}

func (t *execHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query, args, err := decodePayload(r)
	if err != nil {
		ErrorHandler("bad request", err.Error(), http.StatusBadRequest).ServeHTTP(w, r)
		return
	}
	err = t.env.database.Exec(query, args...)
	if err != nil {
		ErrorHandler("bad request", "failed to query data: "+err.Error(), http.StatusBadRequest).ServeHTTP(w, r)
		return
	}
	w.Write([]byte{})

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
	log.Println(t.Error)
	err := json.NewEncoder(w).Encode(map[string]any{
		"error":   t.Error,
		"message": t.Message,
		"code":    t.Code},
	)
	if err != nil {
		http.Error(w, "server error: failed to send error", http.StatusInternalServerError)
	}
}
