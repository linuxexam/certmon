package main

import (
	"log"
	"testing"
)

func TestGetUserCerts(t *testing.T) {
	db := InitDB()
	certs := db.GetUserCerts("user01")
	if len(certs) != 2 {
		t.Fail()
	}
}
func InitDB() *DB {
	testDB, err := NewDB(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	if err := testDB.InsertSampleData(); err != nil {
		log.Fatal(err)
	}
	return testDB
}
