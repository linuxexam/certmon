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

func main() {
	var listen string
	flag.StringVar(&listen, "listen", ":8080", "listen address and port")
	flag.Parse()

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if DEBUG {
			http.FileServer(http.Dir("ui")).ServeHTTP(w, r)
		} else {
			uiFS, _ := fs.Sub(UI, "ui")
			http.FileServerFS(uiFS).ServeHTTP(w, r)
		}
	})

	http.HandleFunc("GET /fetch", func(w http.ResponseWriter, r *http.Request) {
		// get list of certs registered by the current user
		certs := []*Cert{
			{
				Host: "google.ca",
				Port: 443,
			},
			{
				Host: "ww.utoronto.ca",
				Port: 443,
			},
			{
				Host: "idpz.utorauth.utoronto.ca",
				Port: 443,
			},
			{
				Host: "www2.csdn.net",
				Port: 442,
			},
			{
				Host: "www.csdn.net",
				Port: 443,
			},
		}
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
