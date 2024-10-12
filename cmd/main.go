package main

import (
	"github.com/RobinReichert/d-bee-core/internal"
)

func main() {
	database := internal.Database()
	env := internal.Env(database)
}
