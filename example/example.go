package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/miracl/maas-sdk-go"
)

const (
	serviceName = "rpa-example"
	serviceID   = "rpa-example"
	state       = "test-state"
	seesionKey  = "session-key"
)

var (
	clientID     = flag.String("client-id", "", "OIDC Client Id")
	clientSecret = flag.String("client-secret", "", "OIDC Client Secret")
	discovery    = flag.String("discovery", "", "OIDC Discovery URL")
	redirectURL  = flag.String("redirect", "", "Redirect URL")
	addr         = flag.String("addr", ":8002", "Listen address")
	templatesDir = flag.String("templates-dir", "templates", "Template files location")
	debug        = flag.Bool("debug", false, "Debug mode")

	mc maas.Client
)

func checkSession(r *http.Request, sessions map[string]maas.UserInfo) (user maas.UserInfo, err error) {
	c, err := r.Cookie("session")
	if err != nil {
		return user, err
	}
	user, ok := sessions[c.Value]
	if !ok {
		return maas.UserInfo{}, fmt.Errorf("Session %v does not exist", c.Value)
	}
	return user, err
}

func createSession(w http.ResponseWriter, user maas.UserInfo, sessions map[string]maas.UserInfo) {
	sessionID := time.Now().Format(time.RFC850)
	sessions[sessionID] = user
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{Name: "session", Value: sessionID, Expires: expiration}
	http.SetCookie(w, &cookie)
}

func deleteSession(r *http.Request, w http.ResponseWriter, sessions map[string]maas.UserInfo) {
	c, _ := r.Cookie("session")
	delete(sessions, c.Value)
	c.Value = ""
	c.Expires = time.Time{}
	http.SetCookie(w, c)
}

type flash struct {
	Category string
	Message  string
}

type context struct {
	Messages   []flash
	Retry      bool
	AuthURL    string
	Authorized bool
	Email      string
	UserID     string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	flag.Parse()

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

	mc, err := maas.NewClient(maas.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		RedirectURI:  *redirectURL,
		DiscoveryURI: *discovery,
	})
	if err != nil {
		log.Fatal(err)
	}

	sessions := map[string]maas.UserInfo{}

	http.HandleFunc("/oidc", func(w http.ResponseWriter, r *http.Request) {
		ctx := context{}
		ctx.Messages = make([]flash, 10)

		code := r.URL.Query().Get("code")
		accessToken, jwt, err := mc.ValidateAuth(code)
		if err != nil {
			authURL, e := mc.GetAuthRequestURL("test-state")
			if e != nil {
				ctx.Messages = append(ctx.Messages, flash{Category: "error", Message: e.Error()})
			}
			ctx.AuthURL = authURL
		}
		if *debug {
			claims, err := jwt.Claims()
			log.Printf("Access token: %v", accessToken)
			log.Printf("JTW payload: %+v", claims)
		}

		user, err := mc.GetUserUnfo(accessToken)
		if err != nil {
			ctx.Messages = append(ctx.Messages, flash{Category: "error", Message: err.Error()})
		} else {
			createSession(w, user, sessions)
			http.Redirect(w, r, "/", 302)
			return
		}
		if t, err := template.New("index.tmpl").ParseFiles(filepath.Join(*templatesDir, "index.tmpl")); err != nil {
			log.Fatalf("Failed to parse template: %+v", err)
		} else {
			t.Execute(w, ctx)
		}
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		deleteSession(r, w, sessions)
		http.Redirect(w, r, "/", 302)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		ctx := context{}
		ctx.Messages = make([]flash, 0)

		if user, err := checkSession(r, sessions); err != nil {
			authURL, e := mc.GetAuthRequestURL("test-state")
			if e != nil {
				ctx.Messages = append(ctx.Messages, flash{Category: "error", Message: e.Error()})
			}
			ctx.AuthURL = authURL
		} else {
			ctx.Authorized = true
			ctx.UserID = user.UserID
			ctx.Email = user.Email
		}

		if t, err := template.New("index.tmpl").ParseFiles(filepath.Join(*templatesDir, "index.tmpl")); err != nil {
			log.Fatalf("Failed to parse template: %+v", err)
		} else {
			if err = t.Execute(w, ctx); err != nil {
				log.Println(err)
			}
		}
	})

	log.Printf("Service %s started. Listening on %s", serviceName, *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
