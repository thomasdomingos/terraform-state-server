terraform {
  backend "http" {
    address = "http://localhost:8080/state/new"
    lock_address = "http://localhost:8080/state/new"
    unlock_address = "http://localhost:8080/state/new"
  }
}

resource "null_resource" "example" {
  provisioner "local-exec" {
    command = "sleep 30"
  }
}
