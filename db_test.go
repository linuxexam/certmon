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

func TestAddUser(t *testing.T) {
	db := InitDB()
	err := db.AddUser("user02", "user02@gmail.com")
	if err != nil {
		log.Fatal(err)
	}
}

func TestAddUserCert(t *testing.T) {
	db := InitDB()
	err := db.AddUserCert("user01", "2.2.2.2:443", "")
	if err != nil {
		log.Fatal(err)
	}
	if len(db.GetUserCerts("user01")) != 3 {
		t.Fail()
	}
}

func TestDelUserCert(t *testing.T) {
	db := InitDB()
	err := db.DelUserCert("user01", "google.ca:443", "")
	if err != nil {
		log.Fatal(err)
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
