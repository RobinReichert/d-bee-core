package main

import (
	"github.com/RobinReichert/d-bee-core/internal"
)

func main() {
	database := internal.PostgresDatabase()
	env := internal.Env(&database)
	router := internal.Router(env)
	router.Serve()
}
