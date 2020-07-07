package main

import (
	"flag"
	"log"

	"github.com/thomasdomingos/terraform-state-server/config"
	"github.com/thomasdomingos/terraform-state-server/server"
	"github.com/thomasdomingos/terraform-state-server/states"
)

var cfg config.Config

func main() {
	var configPath string

	flag.StringVar(&configPath, "configpath", "/etc/config.yaml", "Configuration file")
	flag.Parse()

	// Load configuration
	if err := config.InitConfig(configPath, &cfg); err != nil {
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
