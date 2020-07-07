// Package server provides the implementation of the services endpoints
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thomasdomingos/terraform-state-server/config"
	"github.com/thomasdomingos/terraform-state-server/states"
	"io/ioutil"
	"log"
	"net/http"
)

var registryPath string
var manager states.Mgr

func homepage(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT: homepage")
	fmt.Fprintf(w, "Welcome to Terraform HTTP State Server!")
}

func getStates(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT: getStates")

	type State struct {
		Name   string   `json:"name"`
		Layers []string `json:"history"`
	}
	States := make([]State, 0)
	for _, name := range manager.GetStates() {
		history := manager.GetHistory(name)
		States = append(States, State{Name: name, Layers: history})
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(States)
}

func getState(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT: getState")
	vars := mux.Vars(r)
	stateName := vars["id"]
	b, err := manager.GetState(stateName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func copy(src, dst string) {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	} // Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func putState(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT: putState")
	vars := mux.Vars(r)
	stateName := vars["id"]

	// Read content of the POST data (verify correctness of json)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !json.Valid(reqBody) {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(errors.New("invalid json data: aborting"))
	}
	// FInally write state
	err = manager.PutState(stateName, reqBody)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Serve configures the routes and start the server using ListenAndServe
func Serve(cfg config.Config, mgr states.Mgr) error {
	log.Println("Serving Terraform HTTP State Server")
	manager = mgr
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homepage)
	router.HandleFunc("/states", getStates).Methods("GET")
	router.HandleFunc("/state/{id}", getState).Methods("GET")
	router.HandleFunc("/state/{id}", putState).Methods("POST")
	registryPath = cfg.Registry.Path
	return http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port), router)
}
