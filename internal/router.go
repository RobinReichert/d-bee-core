package internal

import (
	"net/http"
)

type router struct {
	mux *http.ServeMux
}

func Router(env env) router {
	mux := http.NewServeMux()
	router := router{mux: mux}
	return router

}
