package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/go-chi/chi/v5"
)

var apiKey, apiSecret string
var streamClient *stream.Client

var clientTmpl, jsTmpl *template.Template

func init() {
	var ok bool
	var k string

	k = "GETSTREAMIO_API_KEY"
	apiKey, ok = os.LookupEnv(k)
	if !ok {
		panic(errors.New(k))
	}

	k = "GETSTREAMIO_API_SECRET"
	apiSecret, ok = os.LookupEnv(k)
	if !ok {
		panic(errors.New(k))
	}

	var err error

	clientTmpl, err = template.ParseFiles("client.html")
	if err != nil {
		panic(err)
	}

	jsTmpl, err = template.ParseFiles("client.js")
	if err != nil {
		panic(err)
	}

	streamClient, err = stream.NewClient(apiKey, apiSecret)
	if err != nil {
		panic(err)
	}
}

func main() {
	r := chi.NewRouter()

	r.Get("/users/{uid}/token", getUserToken)
	r.Get("/client{num}.html", getClient)
	r.Get("/client.js", getJS)
	r.Get("/streamchat.js", getStreamchatJS)

	log.Printf("Starting server on port :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func getUserToken(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	t, err := streamClient.CreateToken(uid, time.Time{})
	if err != nil {
		log.Println("Error creating token:", err)
		http.Error(w, "Unexpected error", http.StatusInternalServerError)

		return
	}

	users := []*stream.User{
		{
			ID:   uid,
			Name: uid,
			Role: "admin",
		},
	}

	resp, err := streamClient.UpsertUsers(context.Background(), users...)
	if err != nil {
		panic(err)
	}

	log.Printf("User updated: %+v", resp)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(t))
}

// Returns an HTML page that creates (or connects to) the shared
// channel, creates a new user and starts sending messages on behalf
// of the user.
func getClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	num := chi.URLParam(r, "num")

	err := clientTmpl.Execute(w, struct{ Num string }{num})
	if err != nil {
		log.Println("Error executing client.html template:", err)
		http.Error(w, "Unexpected error", http.StatusInternalServerError)

		return
	}
}

func getJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")

	err := jsTmpl.Execute(w, struct{ APIKey string }{apiKey})
	if err != nil {
		log.Println("Error executing client.js template:", err)
		http.Error(w, "Unexpected error", http.StatusInternalServerError)

		return
	}
}

func getStreamchatJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, "./streamchat.js")
}
