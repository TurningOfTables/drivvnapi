# drivvn backend test

## Description

Basic API setup to the specs of https://github.com/drivvn/backend-task.

Written in Go using the (Go Fiber Framework)[https://gofiber.io/], and using an SQLite database.

## Building and running

Built binaries are in /bin for Windows, Mac and Linux, just run to start the API - defaults to `localhost:8000`

You can build from source with `go build`.

If you provide the `-r` flag the database will be reset to contain one car and the standard four colours at startup.

## Usage

Open API spec for API usage can be found in ./openapi.yml

## Testing

Run `go test` to run tests