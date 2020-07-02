package main

import (
	"log"

	"github.com/thomasdomingos/terraform-state-server/config"
	"github.com/thomasdomingos/terraform-state-server/server"
)

var cfg config.Config

func main() {
	if err := config.InitConfig("./config.yaml", &cfg); err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.Serve(cfg))
}
