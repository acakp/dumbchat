package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/layout.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/chat", chatHandler)
	fmt.Println("starting at :8888...")
	http.ListenAndServe(":8888", r)
}
