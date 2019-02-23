# Trip Booking Service 

This service receives trip requests from a rider client and matches them with a 
willing and available driver

## Setup Development Environment

* Get the suitable `Go` compiler for your machine from `https://golang.org/dl/`
* Get this code `go get -u github.com/pancakem/rides` 
* From `rides/v1/src/cmd/app/` directory run `go get -d ./...`
* Get a redis image by running `docker pull redis`
* Start a containerised redis instance by running this `docker run --name some-redis -d -p 6379 redis`
* Download [Postgres](https://www.postgresql.org/download/) then run the file at `rides/v1/src/pkg/model/schema.sql`
* Change the `config.yaml` file to reflect your system database
* Finally run the server `go run ./v1/src/cmd/app/`
