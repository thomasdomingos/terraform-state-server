package main

import (
	"github.com/thomasdomingos/terraform-state-server/server"
)

func main() {
	server.Serve("8080")
}
