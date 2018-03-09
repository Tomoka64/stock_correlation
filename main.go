package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	router := mux.NewRouter()
	http.Handle("/", router)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "notfound.html", nil)
	})
	tpl = template.Must(template.ParseGlob("templates/*"))
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/result", Result).Methods("GET")
	router.HandleFunc("/ranking", ranking).Methods("GET")
	router.Handle("/favicon.ico", router.NotFoundHandler)
	router.PathPrefix("/publics/").Handler(http.StripPrefix("/publics/", http.FileServer(http.Dir("publics/"))))
}
