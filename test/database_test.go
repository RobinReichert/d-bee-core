package test

import (
	"testing"

	"github.com/RobinReichert/d-bee-core/internal"
)

func TestQuery(t *testing.T) {
	db := internal.PostgresDatabase()
	_, err := db.Query("SELECT * FROM test")
	if err != nil {
		t.Fail()
	}
	_, err = db.Query("SELECT * FROM fail")
	if err == nil {
		t.Fail()
	}
	result, err := db.Query("SELECT * FROM test WHERE id = 1")
	if err != nil {
		t.Fail()
	}
	if len(result) != 1 {
		t.Fail()
	}
	result, err = db.Query("SELECT * FROM test WHERE id = $1", 1)
	if err != nil {
		t.Fail()
	}
	if len(result) != 1 {
		t.Fail()
	}
	result, err = db.Query("SELECT * FROM test WHERE id = 4")
	if err != nil {
		t.Fail()
	}
	if result[0]["email"] != "dave@example.com" {
		t.Fail()
	}
	result, err = db.Query("SELECT status FROM test WHERE id = 4")
	if err != nil {
		t.Fail()
	}
	if result[0]["status"] != false {
		t.Fail()
	}
}
