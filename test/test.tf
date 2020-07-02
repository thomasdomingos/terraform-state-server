terraform {
  backend "http" {
    address = "http://localhost:8080/state"
  }
}

resource "null_resource" "cluster" {
}
