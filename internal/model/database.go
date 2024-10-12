package internal

type dbConnection interface {
	Exec(query string, args ...any) (queryResult, error)
}

type queryResult struct {
}

type database struct {
}

type testDatabase struct {
}
