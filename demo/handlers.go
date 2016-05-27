package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// StatusResponse holds name and version info
type StatusResponse struct {
	ServiceName    string `json:"service_name"`
	ServiceVersion string `json:"service_version"`
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	status := StatusResponse{
		ServiceName:    ServiceName,
		ServiceVersion: ServiceVersion,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {

	authURL, err := mc.GetAuthRequestURL("test-state")
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("unable to get auth url: %v", err))
		return
	}
	log.Printf("Auth URL: %v", authURL)
	t := template.New("index.tmpl")
	t, _ = t.ParseFiles(filepath.Join(*templatesDir, "index.tmpl"))
	t.Execute(w, map[string]string{
		"AuthURL": authURL,
	})
	log.Println("Executed template")
}

func handleOIDCLogin(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "code query param must be set")
		return
	}

	state := r.URL.Query().Get("state")
	if state != "test-state" {
		writeError(w, http.StatusBadRequest, "state query param is invalid")
		return
	}

	_, token, err := mc.ValidateAuth(code)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to validate login: %v", err))
		return
	}
	claims, err := token.Claims()
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("fauiled to retrieve claims: %v", err))
		return
	}

	user, ok, err := claims.StringClaim("sub")

	var s string
	if err != nil {
		s = fmt.Sprintf("No sub. %v", err)
	} else if !ok {
		s = fmt.Sprintf("No sub")
	} else {
		s = fmt.Sprintf("<html>Logged in as <b>%s</b></br><a href='/'>Logout</a></html>", user)
	}

	w.Write([]byte(s))
}
