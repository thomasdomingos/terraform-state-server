// Package server provides the implementation of the services endpoints
package server

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"

  "github.com/gorilla/mux"
)

func homepage(w http.ResponseWriter, r *http.Request) {
  log.Println("HIT: homepage")
  fmt.Fprintf(w, "Welcome to Terraform HTTP State Server!")
}

func assertFile(name string) error {
  _, err := os.Stat(name)
  if err == nil {
    return nil
  }
  if os.IsNotExist(err) {
    log.Println("File does not exist, creating it")
    file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
    defer file.Close()
    if err != nil {
      return err
    }
  }
  return nil
}

func getState(w http.ResponseWriter, r *http.Request) {
  log.Println("HIT: getState")
  if err := assertFile("/tmp/state"); err != nil {
    log.Fatal(err)
  }
  file, err := os.Open("/tmp/state")
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

func postState(w http.ResponseWriter, r *http.Request) {
  log.Println("HIT: postState")
  if err := assertFile("/tmp/state"); err != nil {
    log.Fatal(err)
  }

  reqBody, err := ioutil.ReadAll(r.Body)

  file, err := os.OpenFile("/tmp/state", os.O_WRONLY, 0644)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  n, err := file.Write(reqBody)
  if err != nil {
    log.Fatal(err)
  }
  log.Println(n, "bytes written to file")
}

// Serve configures the routes and start the server using ListenAndServe
func Serve(port string) {
  log.Println("Serving Terraform HTTP State Server on port", port)
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", homepage)
  router.HandleFunc("/state", getState).Methods("GET")
  router.HandleFunc("/state", postState).Methods("POST")
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
