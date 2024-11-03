# Go Cash
Go Cash is a fictional centralized economy, built for the purpose of learning Go.

## Requirements

Go Cash is not meant to be directly exposed to the internet. You should configure a webserver such as Nginx or Caddy, to proxy the application. You may also want to consider configuring it to behave as a service, depending on your environment.

`go 1.23.2`

## Installation

### Setup the database

`git clone https://github.com/Hypnophobe/go-cash && cd go-cash && go run initdb/main.go`

### Build the binary and start the server

```
go build .
./go-cash
```