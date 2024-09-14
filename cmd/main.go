package main

import (
	"NotesAPI/handlers"
	"NotesAPI/internal"
	"NotesAPI/utils"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	log.Print("Prepare db...")
	if err := internal.Prepare(); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/public/login", handlers.Login)

	http.Handle("/protected/addnote", utils.AuthenticationMiddleware(http.HandlerFunc(handlers.AddNotesHandler)))
	http.Handle("/protected/getnotes", utils.AuthenticationMiddleware(http.HandlerFunc(handlers.GetNotesHandler)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
