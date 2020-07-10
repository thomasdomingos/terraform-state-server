# terraform-state-server

Simple File server to be plugged with terraform http state

## Installation

### Run as standalone application
- checkout the source and `cd` into the checked out directory
- change `config.yaml` and Dockerfile according to your needs (see [configuration](## Configuration))
- build frm the source : `go build -o <name of your executable> main.go`
- run it : `./<name of your executable>`


### Run in Docker

- checkout the source and `cd` into the checked out directory
- change `config.yaml` and Dockerfile according to your needs (see [configuration](## Configuration))
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

## API

The server uses a simple RestAPI. The following presume it is running on localhost:8080 :

`GET http://localhost:8080/state/<my_state>`

Get last version of the state named `<my_state>`.

`POST http://localhost:8080/state/<my state>`

Update the last version of the state named `<my_state>`.

`LOCK http://localhost:8080/state/my_state`

Put a lock on the state `<my_state>`, return OK (200) if the lock is acquired, StatusLock (423) is the lock is already taken.

`UNLOCK http://localhost:8080/state/my_state`

Unlock state `<my_state>`.

Get last version of the state named `<my_state>`

`GET http://localhost:8080/states`

Get all states present in the server, along with the identifiers of all the version of the state.

## Configuration

Configuration is equivalently done from file or environment variables. The file at `/etc/config.yaml` is used by default but other configuration file can be specified, using the parameter -config.

```
[server]
host="localhost" # or env SERVER_HOST
port="8080" # or env SERVER_PORT

[registry]
path="./tss" # or env REGISTRY_PATH

[database]
path="./tss.db" # or env DATABASE_PATH
```

## Internal

### States
States are preserved as files, named after their fingerprint. This fingerprint is computed from the state name, the content of the state, and the fingerprint of the previous state. This will allow to run checks on the integrity of the data in a future release.
These files are kept under the directory `regostry` in the configuration file

### Database
A database keeps track of the files that correspond to each state.

## Future work :
- Unit tests
- Simple user management to enable BasicAuth.
- Simple history management to recover a previous state.

## Licence
[GNU General Public License v3.0](https://github.com/thomasdomingos/terraform-state-server/blob/master/LICENSE)
