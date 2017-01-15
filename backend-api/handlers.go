package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/mangirdaz/ocp-demo/config"
)

//keep order as it is as we use it in json unmarshal. if needed add to bottom
type Note struct {
	Note         string `json:"note"`
	Key          string `json:"key"`
	Path         string `json:"path"`
	CreationTime string `json:"creation_time"`
}

type Notes struct {
	Note []Note
}

//Index index method for API
func Index(w http.ResponseWriter, r *http.Request, storage *LibKVBackend) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OCP Nodes Demo API")
}

//extend standart handler with our required storage backend details
type backendHandler func(w http.ResponseWriter, r *http.Request, storage *LibKVBackend)

type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

//return what mux expects
func mybackendHandler(handler backendHandler, storage *LibKVBackend) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, storage)
	}
}

//Health endpoit for app server
func Health(w http.ResponseWriter, r *http.Request, storage *LibKVBackend) {
	log.Debug("/health endpoint called")
	w.Write([]byte("OK"))
}

func CheckAuth(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	log.Info("Check auth Middleware")
	user, pass, _ := req.BasicAuth()
	enabled, _ := strconv.ParseBool(config.Get("EnvBasicAuth"))
	log.WithFields(log.Fields{
		"auth":  enabled,
		"auth1": config.Get("EnvBasicAuth"),
	}).Debug("handler")
	if enabled && !checkPass(user, pass) {
		reason := "Unauthorized"
		resp.WriteHeader(http.StatusUnauthorized)
		response(reason, true, nil, resp, req)
		return
	}
	next(resp, req)
}

func CorsHeadersMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Info("Cors Middleware")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	rw.Header().Set("Access-Control-Expose-Headers", "Authorization")
	rw.Header().Set("Access-Control-Request-Headers", "Authorization")

	if r.Method == "OPTIONS" {
		rw.WriteHeader(200)
		return
	}

	next(rw, r)
}

func checkPass(user, pass string) bool {
	log.Info(fmt.Sprintf("User [%s] and Pass [%s]", user, pass))
	if user == "admin" && pass == "admin" {
		log.Info("Pass OK")
		return true
	} else {
		log.Info("Pass Error")
		return false
	}
	return false
}

//GetNotes returns all notes
func GetNotes(resp http.ResponseWriter, req *http.Request, storage *LibKVBackend) {
	var note Note
	note.Path = "notes"

	//get all notes
	log.Info("Get notes")
	notes, err := storage.GetAll(note.Path)
	if err != nil {
		log.Error(err)
	}

	//create json on all notes
	b, err := json.Marshal(notes)
	if err != nil {
		log.Errorf("Error: %s", err)
	}

	resp.WriteHeader(http.StatusOK)
	result := string(b)
	response(result, true, nil, resp, req)
}

//GetNotes returns all notes
func CreateNote(resp http.ResponseWriter, req *http.Request, storage *LibKVBackend) {
	var note Note
	note.Path = "notes"

	time := time.Now()

	//decode json
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&note)
	if err != nil {
		log.Error(err)
	}
	defer req.Body.Close()
	note.Key = config.GenerateID()
	//parse notes to object and validate

	note.CreationTime = time.String()

	b, err := json.Marshal(note)
	if err != nil {
		fmt.Println("error:", err)
	}

	b, err = json.Marshal(note)
	if err != nil {
		log.Error(err)
	}
	err = storage.Put(note.Key, b)
	if err != nil {
		log.Error(err)
	}

	log.Info(fmt.Sprintf("Create note [%s]", note))
}

//GetNote get one note
func GetNote(w http.ResponseWriter, r *http.Request, storage *LibKVBackend) {

	vars := mux.Vars(r)
	key := vars["key"]
	var note Note
	note.Path = "notes"

	//validate key(id) of the node and check if it exist
	log.Info(fmt.Sprintf("Get note [%s/%s]", note.Path, key))

}

//DeleteNote delete one node
func DeleteNote(w http.ResponseWriter, r *http.Request, storage *LibKVBackend) {

	vars := mux.Vars(r)
	key := vars["key"]
	var note Note
	note.Path = "notes"

	//validate key(id) of the node and check if it exist
	log.Info(fmt.Sprintf("Delete note [%s/%s]", note.Path, key))

}

//UpdateNote update one node
func UpdateNote(w http.ResponseWriter, r *http.Request, storage *LibKVBackend) {

	vars := mux.Vars(r)
	key := vars["key"]
	var note Note
	note.Path = "notes"

	//validate key(id) of the node and check if it exist
	log.Info(fmt.Sprintf("Update note [%s/%s]", note.Path, key))

}

//GetExternal get external content
func GetExternal(w http.ResponseWriter, r *http.Request, storage *LibKVBackend) {

	vars := mux.Vars(r)
	key := vars["key"]
	var note Note
	note.Path = "notes"

	//validate key(id) of the node and check if it exist
	log.Info(fmt.Sprintf("Get external content [%s/%s]", note.Path, key))

}

func response(obj interface{}, prettyPrint bool, err error, resp http.ResponseWriter, req *http.Request) {
	// Check for an error
HAS_ERR:
	if err != nil {

		if err.Error() == NotFound {
			resp.WriteHeader(http.StatusNotFound)
			return
		}

		log.WithFields(log.Fields{
			"error":  err,
			"method": req.Method,
			"url":    req.URL,
		}).Error("request error")

		code := 500
		errMsg := err.Error()
		if strings.Contains(errMsg, "Permission denied") || strings.Contains(errMsg, "ACL not found") {
			code = 403
		}
		resp.WriteHeader(code)
		resp.Write([]byte(err.Error()))
		return
	}

	// Write out the JSON object
	if obj != nil {
		buf, err := marshall(obj, true)
		if err != nil {
			goto HAS_ERR
		}
		resp.Header().Set("Content-Type", "application/json")

		// encoding/json library has a specific bug(feature) to turn empty slices into json null object,
		// let's make an empty array instead
		if string(buf) == "null" {
			buf = []byte("[]")
		}
		resp.Write(buf)
	}
}

// marshall returns a json byte slice, leaving existing json untouched.
func marshall(obj interface{}, pretty bool) ([]byte, error) {

	var js interface{}
	var buf []byte

	// Only check objects that byte slices and strings for valid json
	switch v := obj.(type) {
	case []byte:
		buf = []byte(v)
	case string:
		buf = []byte(v)
	}

	// If we were given a valid json object, return it as-is
	if buf != nil && json.Unmarshal(buf, &js) == nil {
		return buf, nil
	}

	// Otherwise marshall the object into json
	if pretty {
		return json.MarshalIndent(obj, "", "    ")
	}
	return json.Marshal(obj)
}
