package main

/*
* Database
 */

type DB struct {
	connStr string
}

func (db *DB) GetUserCerts(userId string) []*Cert {
	return nil
}

func (db *DB) GetAllCerts() []*Cert {
	return nil
}

func (db *DB) AddCert(host string, port int) {

}

func (db *DB) RemoveCert(host string, port int) {

}

func (db *DB) AddOwner(host string, port int, userId string) {

}

func (db *DB) RemoveOwner(host string, port int, userId string) {
}
