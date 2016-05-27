package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/miracl/maas-sdk-go"
)

var (
	clientID     = flag.String("client-id", "", "OIDC Client Id")
	clientSecret = flag.String("client-secret", "", "OIDC Client Secret")
	discovery    = flag.String("discovery", "", "OIDC Discovery URL")
	redirectURL  = flag.String("redirect", "", "Redirect URL")
	printVersion = flag.Bool("version", false, "Print the program version")
	addr         = flag.String("addr", ":8002", "Listen address")
	staticDir    = flag.String("static-dir", "static", "Static files location")
	codepadDir   = flag.String("codepad-dir", "codepad", "Codepad files location")
	templatesDir = flag.String("templates-dir", "templates", "Template files location")

	mc maas.Client
)

func initClient() maas.Client {
	c, err := maas.NewClient(maas.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		RedirectURI:  *redirectURL,
		DiscoveryURI: *discovery,
	})
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func writeError(w http.ResponseWriter, code int, msg string) {
	e := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}
	b, err := json.Marshal(e)
	if err != nil {
		log.Printf("Failed marshaling %#v to JSON: %v", e, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	flag.Parse()

	if *printVersion {
		fmt.Println(ServiceVersion)
		return
	}

	if *clientID == "" {
		log.Fatal("client-id required")
	}

	if *clientSecret == "" {
		log.Fatal("client-secret required")
	}

	if *discovery == "" {
		log.Fatal("Discovery URL required")
	}

	if *redirectURL == "" {
		log.Fatal("Redirect URL required")
	}

	mc = initClient()

	http.Handle("/mpin/", http.StripPrefix("/mpin/", http.FileServer(http.Dir(*codepadDir))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*staticDir))))
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/oidc", handleOIDCLogin)
	http.HandleFunc("/", handleIndex)

	log.Printf("Service %s v %s started. Listening on %s", ServiceName, ServiceVersion, *addr)
	http.ListenAndServe(*addr, nil)
}
