package main

import (
	"crypto/x509"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/linuxexam/certmon"
	"github.com/linuxexam/certmon/data"
)

//go:embed ui
var UI embed.FS
var DEBUG = false

func main() {
	var listen string
	var dataDir string
	var rootCAFile string

	flag.StringVar(&listen, "listen", ":8080", "listen address and port")
	flag.StringVar(&dataDir, "dataDir", ".", "directory for saving data")
	flag.StringVar(&rootCAFile, "rootCAs", "", "PEM file containing root CAs")

	flag.Parse()

	var rootCAs *x509.CertPool

	if rootCAFile != "" {
		rootCAs = x509.NewCertPool()
		pemBytes, err := os.ReadFile(rootCAFile)
		if err != nil {
			log.Fatal(err)
		}
		if ok := rootCAs.AppendCertsFromPEM(pemBytes); !ok {
			panic("failed to parse root certs")
		}
	}

	// db
	db, err := data.NewDB(strings.TrimSuffix(dataDir, "/") + "/certmon.sqlite")
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

	http.HandleFunc("/userid", func(w http.ResponseWriter, r *http.Request) {
		userId := GetCurrentUserId(r)
		fmt.Fprintf(w, "%s", userId)
	})

	// add a cert for a user
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		userId := GetCurrentUserId(r)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		certAddr := strings.Trim(r.FormValue("certAddr"), " \t\r\n")
		certDNS := strings.Trim(r.FormValue("certDNS"), " \t\r\n")

		log.Print(r.Form)

		err = db.AddUserCert(userId, certAddr, certDNS)
		if err != nil {
			log.Print(err.Error())
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
		userId := GetCurrentUserId(r)
		idUserCert, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			fmt.Fprintf(w, "Error:%s", err.Error())
			log.Print(err.Error())
			return
		}

		err = db.DelUserCertById(idUserCert, userId)
		if err != nil {
			fmt.Fprintf(w, "Error:%s", err.Error())
			log.Print(err.Error())
			return
		}
	})

	// get all cert for a user
	http.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		// get list of certs registered by the current user
		certs := db.GetUserCerts(GetCurrentUserId(r))
		// connect to verify each cert
		var wg sync.WaitGroup
		wg.Add(len(certs))
		for _, cert := range certs {
			go func(cert *certmon.Cert) {
				defer wg.Done()
				log.Printf("checking cert at %s(%s)...", cert.Addr, cert.DNS)
				cert.Update(rootCAs)
			}(cert)
		}
		wg.Wait()

		// sort certs by daysLeft
		slices.SortFunc(certs, func(a, b *certmon.Cert) int {
			return a.DaysLeft - b.DaysLeft
		})

		buf, err := json.Marshal(certs)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if _, err := w.Write(buf); err != nil {
			log.Print(err.Error())
		}
	})

	log.Fatal(http.ListenAndServe(listen, nil))
}
