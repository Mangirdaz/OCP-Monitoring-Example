package main

import (
	"net/http"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mangirdaz/ocp-mon-demo/config"
)

func NewRouter() {

	//create new router
	router := mux.NewRouter().StrictSlash(false)
	storage := InitKVStorage()

	//api backend init
	router.Path("/healthz").Name("Health endpoint").HandlerFunc(http.HandlerFunc(mybackendHandler(Health, storage)))
	//api backend init
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	apiV1.Methods("GET").Path("").Name("Index").Handler(http.HandlerFunc(mybackendHandler(Index, storage)))

	//notes methods
	notes := apiV1.PathPrefix("/notes").Subrouter()

	notes.Methods("GET").Name("Notes").Handler(http.HandlerFunc(mybackendHandler(GetNotes, storage)))
	notes.Methods("POST").Name("Create Note").Handler(http.HandlerFunc(mybackendHandler(CreateNote, storage)))

	notes.Methods("DELETE").Path("/{1}").Name("Delete Single Note").Handler(http.HandlerFunc(mybackendHandler(DeleteNote, storage)))
	notes.Methods("GET").Path("/{1}").Name("Get Single Note").Handler(http.HandlerFunc(mybackendHandler(GetNote, storage)))
	notes.Methods("PUT").Path("/{1}").Name("Update Single Note").Handler(http.HandlerFunc(mybackendHandler(UpdateNote, storage)))

	//external resources
	static := apiV1.PathPrefix("/static").Subrouter()
	static.Methods("GET").Name("Index Static").Handler(http.HandlerFunc(mybackendHandler(Index, storage)))
	static.Methods("GET").Path("/{key}").Name("static content").Handler(http.HandlerFunc(mybackendHandler(GetExternal, storage)))

	//middleware intercept
	midd := http.NewServeMux()
	midd.Handle("/", router)
	midd.Handle("/api/v1/notes", negroni.New(
		negroni.HandlerFunc(CorsHeadersMiddleware),
		negroni.Wrap(apiV1),
	))
	n := negroni.Classic()
	n.UseHandler(midd)
	url := fmt.Sprintf("%s:%s", config.Get("EnvAPIIP"), config.Get("EnvAPIPort"))

	log.WithFields(log.Fields{
		"url": url,
	}).Debug("api: starting api server")

	log.Fatal(http.ListenAndServe(url, n))

	//return router

}
