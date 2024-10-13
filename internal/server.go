package internal

import "net/http"

func Server(router router) {
	err := http.ListenAndServe(":80", router.mux)
	if err != nil {
		panic(err)
	}
}
