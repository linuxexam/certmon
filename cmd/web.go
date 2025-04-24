package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/linuxexam/certmon"
	"github.com/linuxexam/certmon/data"
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
	db, err := data.NewDB("./certmon.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// web
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if DEBUG {
			http.FileServer(http.Dir("ui")).ServeHTTP(w, r)
		} else {
			uiFS, _ := fs.Sub(UI, "ui")
			http.FileServerFS(uiFS).ServeHTTP(w, r)
		}
	})

	// add a cert for a user
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		userId := GetCurrentUser()
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		certAddr := r.FormValue("certAddr")
		certDNS := r.FormValue("certDNS")

		log.Print(r.Form)

		err = db.AddUserCert(userId, certAddr, certDNS)
		if err != nil {
			log.Print(err)
			w.Write([]byte(err.Error()))
		}
		w.Write([]byte("good"))
	})

	// delete a cert for a user
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		userId := GetCurrentUser()
		idUserCert, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			fmt.Fprintf(w, "Error:%s", err.Error())
			log.Print(err)
			return
		}

		err = db.DelUserCertById(idUserCert, userId)
		if err != nil {
			fmt.Fprintf(w, "Error:%s", err.Error())
			log.Print(err)
			return
		}
	})

	// get all cert for a user
	http.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		// get list of certs registered by the current user
		certs := db.GetUserCerts(GetCurrentUser())
		// connect to verify each cert
		var wg sync.WaitGroup
		wg.Add(len(certs))
		for _, cert := range certs {
			go func(cert *certmon.Cert) {
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
