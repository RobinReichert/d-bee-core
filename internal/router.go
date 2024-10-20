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

	mux.Handle("/query", QueryHandler(&env))
	mux.Handle("/exec", ExecHandler(&env))

	return router
}

func (t *router) Serve() {
	err := http.ListenAndServe(":80", t.mux)
	if err != nil {
		panic(err)
	}
}
