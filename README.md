# sniplate

`sniplate` is a _template_ for building secure, scalable, APIs with Go. The
_template_ uses text _snippets_ as the initial model to get started. There is
also a user model, you can easily decouple these into micro services as well!

The initial project set up has been completed, and includes the following
features so far.

- 100% Go standard library
- CLI flags for custom server config
- Structured logging
- Centralized error handling
- Middleware
- Security enhancements

## Routes

| Method | URL Pattern     | Handler            | Action                       |
| ------ | --------------- | ------------------ | ---------------------------- |
| GET    | /v1/healthcheck | healthcheckHandler | Show application information |
| POST   | /v1/snips       | createSnipHandler  | Add snip                     |
| GET    | /v1/snips/{id}  | showSnipHandler    | Show specific snip           |

## Getting Started

To get started locally, make sure you have Git and Go installed, then pull the
repository.

### Pull Project

```bash
git pull https://github.com/pwilliams-ck/sniplate
```

Or use Go `get`.

```bash
go get https://github.com/pwilliams-ck/sniplate
```

### Run Project

```bash
go run ./cmd/api
```

You should now be able to run `curl` commands against `localhost:4200`.

```bash
curl -i localhost:4200/v1/healthcheck
```

For server configuration info, try running `go run ./cmd/api -help`.

### Set Up TLS

There is an option to run the server with TLS enabled, you if you need to create
a development TLS certificate and key. If you have Go installed they include a
handy tool to create self-signed certificates with `crypto/tls` package. `cd`
into the project root, `mkdir`, and `cd` again.

```bash
cd sniplate
mkdir tls
cd tls
```

Next, find the `generate_cert.go` tool's path, and run it from the `tls`
directory. Here is the path for MacOS, and probably Linux. The CLI flags are
there for convenience, copy/pasta away.

```bash
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

Then `cd` back into the root directory, build, and run the application.

```bash
cd ..
go run ./cmd/api -tls=true
```

### Build for Remote Server

Remote server testing is completed via the following workflow. This builds the
app for Linux from your local machine, then transfers the executable to the
server.

First, make code changes.

Then run the following, where `./test-api` is the name of the executable you are
building.

```bash
GOOS=linux GOARCH=amd64 go build -o ./test-api ./cmd/api
```

Then run the following to _ssh file transfer_ to move the binary over. Change
the **executable**, **URL**, and **user** to fit your needs.

```bash
scp test-api pwilliams@svc-hub-dev.cloudkey.io:/home/pwilliams/
```

Then `ssh` into the server and run `./api-test`, add `-help` flag for more info.
You might need to make it executable with `chmod +x api-test`.

## Conclusion

This is a work in progress and will be updated regularly.
