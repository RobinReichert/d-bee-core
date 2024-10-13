package internal

type env struct {
	database Database
}

func Env(database Database) env {
	return env{database: database}
}

func (t *env) Database() Database {
	return t.database
}
