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
}
