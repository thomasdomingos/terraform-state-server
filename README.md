# terraform-state-server

Simple File server to be plugged with terraform http state

## Installation

### Run as standalone application
- checkout the source and `cd` into the checked out directory
- change `config.yaml` and Dockerfile according to your needs (see .)
- build frm the source : `go build -o <name of your executable> main.go`
- run it : `./<name of your executable>`


### Run in Docker

- checkout the source and `cd` into the checked out directory
- change `config.yaml` and Dockerfile according to your needs (see .)
- build the docker image : `docker build -t <name_of_your_image> .`
- run it : `docker run -p 8080:8080 <name of your image>`

## Use with terraform
The following terraform file shows how to use your locally running terraform-state-server:
```
terraform {
  backend "http" {
    address = "http://localhost:8080/state/my_state"
    lock_address = "http://localhost:8080/state/my_state"
    unlock_address = "http://localhost:8080/state/my_state"
  }
}

resource "null_resource" "example" {
  provisioner "local-exec" {
    command = "sleep 30"
  }
}
```

# API

The server uses a simple RestAPI. The following presume it is running on localhost:8080 :

`GET http://localhost:8080/state/my_state`

`POST http://localhost:8080/state/my_state`

`LOCK http://localhost:8080/state/my_state`

`UNLOCK http://localhost:8080/state/my_state`

`GET http://localhost:8080/states`

# Future work :
- Tests
- Simple user management to enable BasicAuth

## Licence
[GNU General Public License v3.0](https://github.com/thomasdomingos/terraform-state-server/blob/master/LICENSE)
