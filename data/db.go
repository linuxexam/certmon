package data

import (
	"database/sql"
	"log"

	"github.com/linuxexam/certmon"
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
	_, err = _db.Exec("PRAGMA foreign_keys = ON")
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
				id INTEGER PRIMARY KEY,
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

func (db *DB) GetUserCerts(userId string) []*certmon.Cert {
	sql := `
	SELECT user_cert.id, cert.addr, cert.dns 
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
	var certs []*certmon.Cert
	for rows.Next() {
		var id int
		var addr string
		var dns string
		err := rows.Scan(&id, &addr, &dns)
		if err != nil {
			log.Print(err)
			break
		}
		certs = append(certs, &certmon.Cert{ID: id, Addr: addr, DNS: dns})
	}

	return certs
}

func (db *DB) GetAllCerts() []*certmon.Cert {
	return nil
}

func (db *DB) UpdateUser(userId, email string) error {
	_, err := db.Exec(`INSERT INTO user (id, email)
				VALUES (?,?)
				ON CONFLICT(id) DO UPDATE SET
					email = ?
	`, userId, email, email)
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) AddUserIfNotExits(userId string) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO user (id)
				VALUES (?)`, userId)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) AddCertIfNotExits(addr, dns string) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO cert (addr, dns)
				VALUES (?,?)`,
		addr, dns)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) AddUserCert(userId, certAddr, certDNS string) error {
	log.Printf("%s,%s,%s", userId, certAddr, certDNS)
	if err := db.AddUserIfNotExits(userId); err != nil {
		return err
	}
	if err := db.AddCertIfNotExits(certAddr, certDNS); err != nil {
		return err
	}
	sql := `
		INSERT INTO user_cert (user_id, cert_id)
		VALUES (?,(SELECT id FROM cert WHERE addr = ? and dns = ?))
	`
	_, err := db.Exec(sql, userId, certAddr, certDNS)
	return err
}

// only allow deleting certs owned by the userId
func (db *DB) DelUserCertById(id int, userId string) error {
	sql := `DELETE FROM user_cert where id=? and user_id=?`
	_, err := db.Exec(sql, id, userId)
	return err
}
