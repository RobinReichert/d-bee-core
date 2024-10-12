package internal

type env struct {
	database dbConnection
}

func Env(database dbConnection) env {
	return env{database: database}
}
