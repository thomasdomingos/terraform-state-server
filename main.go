package main

import (
	"log"

	"github.com/thomasdomingos/terraform-state-server/config"
	"github.com/thomasdomingos/terraform-state-server/server"
	"github.com/thomasdomingos/terraform-state-server/states"
)

var cfg config.Config

func main() {
	// Load configuration
	if err := config.InitConfig("./config.yaml", &cfg); err != nil {
		log.Fatal(err)
	}

	mgr := states.Mgr{}
	err := mgr.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer mgr.Close()
	log.Fatal(server.Serve(cfg, mgr))

	/*
	  log.Printf("%#v\n", state)
	  newstate := database.NextState(*state, "tutu")
	  log.Printf("%#v\n", newstate)
	*/

	/*  state := states.NewState("toto", "titi")
	    states.InsertState(db, *state)
	    newstate := states.NextState(*state, "tutu")
	    states.InsertState(db, *newstate)
	    if cc, err := states.Get(db, "toto"); err == nil {
	      log.Println(cc)
	    }
	    log.Println(states.GetAll(db))
	*/
}
