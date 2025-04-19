package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

/*
* Database
 */

type DB struct {
	*sql.DB
	Dsn string
}

func NewDB(dsn string) (*DB, error) {
	db := &DB{Dsn: dsn}

	_db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.DB = _db
	// initialize db if not yet
	sql := `
		CREATE TABLE IF NOT EXISTS user(
				id TEXT PRIMARY KEY,
				email TEXT
		);

		CREATE TABLE IF NOT EXISTS cert(
				id INTEGER PRIMARY KEY,
				addr TEXT NOT NULL,
				dns TEXT,
				UNIQUE(addr, dns)
		);

		CREATE TABLE IF NOT EXISTS user_cert(
				user_id TEXT NOT NULL,
				cert_id INTEGER NOT NULL,
				UNIQUE(user_id, cert_id),
				FOREIGN KEY (user_id) REFERENCES user(id),
				FOREIGN KEY (cert_id) REFERENCES cert(id)
		);
	`
	if _, err := db.Exec(sql); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) InsertSampleData() error {
	sql := `
	INSERT INTO cert(addr, dns)
	VALUES
		("google.com:443", NULL),
		("baidu.com:443", NULL),
		("1.2.3.4:443", "myexample.com");
	
	INSERT INTO user(id, email)
	VALUES
		("user01", "user01@gmail.com");

	INSERT INTO user_cert(user_id, cert_id)
	VALUES
		("user01", 1),
		("user01", 3);
	`
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetUserCerts(userId string) []*Cert {
	sql := `
	SELECT cert.addr
	FROM user_cert JOIN cert ON user_cert.cert_id = cert.id
	WHERE user_id = ?
	`
	st, err := db.Prepare(sql)
	if err != nil {
		log.Print(err)
		return nil
	}
	rows, err := st.Query(userId)
	if err != nil {
		log.Print(err)
		return nil
	}
	var certs []*Cert
	for rows.Next() {
		var addr string
		err := rows.Scan(&addr)
		if err != nil {
			log.Print(err)
			break
		}
		certs = append(certs, &Cert{Addr: addr})
	}

	return certs
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
