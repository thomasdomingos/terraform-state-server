terraform {
  backend "http" {
    address = "http://localhost:8080/state/test"
  }
}

resource "null_resource" "example" {
}
