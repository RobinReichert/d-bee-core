package test

import (
	"testing"

	"github.com/RobinReichert/d-bee-core/internal"
)

func TestQuery(t *testing.T) {
	db := internal.PostgresDatabase()
	result, err := db.Query("SELECT * FROM test")
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
