terraform {
  backend "http" {
    address = "http://localhost:8080/state/new"
  }
}

resource "null_resource" "example" {
}
