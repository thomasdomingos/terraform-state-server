// Package server provides the implementation of the services endpoints
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thomasdomingos/terraform-state-server/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var registryPath string

func homepage(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT: homepage")
	fmt.Fprintf(w, "Welcome to Terraform HTTP State Server!")
}

func assertState(stateName string) error {
	// Test directory containing state
	if _, err := os.Stat(filepath.Join(registryPath, stateName)); err != nil {
		// Create directory that will contain state file(s)
		if os.IsNotExist(err) {
			log.Println("Directory does not not exist, creating it")
			if err := os.MkdirAll(filepath.Join(registryPath, stateName), os.ModePerm); err != nil {
				return err
			}
		}
	}
	// Create an empty file when state does not exists
	if _, err := os.Stat(filepath.Join(registryPath, stateName, "state")); err != nil {
		if os.IsNotExist(err) {
			log.Println("Creating empty state file")
			file, err := os.OpenFile(filepath.Join(registryPath, stateName, "state"), os.O_RDONLY|os.O_CREATE, 0644)
			defer file.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getStates(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT: getStates")

	states, err := ioutil.ReadDir(registryPath)
	if err != nil {
		log.Fatal(err)
	}
	type State struct {
		Name string `json:"Name"`
	}
	States := make([]State, 0)
	for _, f := range states {
		if f.IsDir() {
			States = append(States, State{f.Name()})
		}
	}
	json.NewEncoder(w).Encode(States)
}

func getState(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT: getState")
	vars := mux.Vars(r)
	stateName := vars["id"]
	if err := assertState(stateName); err != nil {
		log.Fatal(err)
	}
	// Read state file and return its content
	file, err := os.Open(filepath.Join(registryPath, stateName, "state"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(b)
	if err != nil {
		log.Fatal(err)
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

func postState(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT: postState")
	vars := mux.Vars(r)
	stateName := vars["id"]
	if err := assertState(stateName); err != nil {
		log.Fatal(err)
	}

	// Copy current state to keep history
	copy(filepath.Join(registryPath, stateName, "state"), filepath.Join(registryPath, stateName, strconv.FormatInt(time.Now().Unix(), 10)))

	// Read content of the POST data, and write it to the corresponding state file
	reqBody, err := ioutil.ReadAll(r.Body)
	if !json.Valid(reqBody) {
		log.Fatal(errors.New("invalid json data: aborting"))
	}
	err = ioutil.WriteFile(filepath.Join(registryPath, stateName, "state"), reqBody, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Serve configures the routes and start the server using ListenAndServe
func Serve(cfg config.Config) error {
	log.Println("Serving Terraform HTTP State Server")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homepage)
	router.HandleFunc("/states", getStates).Methods("GET")
	router.HandleFunc("/state/{id}", getState).Methods("GET")
	router.HandleFunc("/state/{id}", postState).Methods("POST")
	registryPath = cfg.Registry.Path
	return http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port), router)
}
