package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"sync"
)

//go:embed ui
var UI embed.FS
var DEBUG = true

func GetCurrentUser() string {
	return "user01"
}

func main() {
	var listen string
	flag.StringVar(&listen, "listen", ":8080", "listen address and port")
	flag.Parse()

	// db
	db, err := NewDB("./certmon.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// web
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if DEBUG {
			http.FileServer(http.Dir("ui")).ServeHTTP(w, r)
		} else {
			uiFS, _ := fs.Sub(UI, "ui")
			http.FileServerFS(uiFS).ServeHTTP(w, r)
		}
	})

	// add a cert for a user
	http.HandleFunc("GET /add", func(w http.ResponseWriter, r *http.Request) {
		userId := GetCurrentUser()
		certAddr := r.URL.Query().Get("certAddr")
		certDNS := r.URL.Query().Get("certDNS")
		err := db.AddUserCert(userId, certAddr, certDNS)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write([]byte("good"))
	})

	// delete a cert for a user
	http.HandleFunc("GET /delete", func(w http.ResponseWriter, r *http.Request) {

	})

	// get all cert for a user
	http.HandleFunc("GET /fetch", func(w http.ResponseWriter, r *http.Request) {
		// get list of certs registered by the current user
		certs := db.GetUserCerts(GetCurrentUser())
		// connect to verify each cert
		var wg sync.WaitGroup
		wg.Add(len(certs))
		for _, cert := range certs {
			go func(cert *Cert) {
				defer wg.Done()
				cert.Update()
			}(cert)
		}
		wg.Wait()

		buf, err := json.Marshal(certs)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if _, err := w.Write(buf); err != nil {
			log.Print(err)
		}
	})

	log.Fatal(http.ListenAndServe(listen, nil))
}
