# sniplate

`sniplate` is for building secure, scalable, HTTP services with Go. This
_template_ uses text snippets (_snips_) as the initial API model to get started,
hence the name _sniplate_. There will also be a user model, you can de-couple
these and build out micro services as well, but I feel this is a good starting
point as a modular monolith.

The initial project set up has been completed, and includes the following
features so far.

- 100% Go standard library
- CLI flags for custom server config
- Structured logging
- Centralized error handling
- Middleware
- Security enhancements
- Docker infrastructure

## Routes

| Method | URL Pattern     | Handler            | Action                       |
| ------ | --------------- | ------------------ | ---------------------------- |
| GET    | /v1/healthcheck | healthcheckHandler | Show application information |
| POST   | /v1/snips       | createSnipHandler  | Add snip                     |
| GET    | /v1/snips/{id}  | showSnipHandler    | Show specific snip           |

## Getting Started

To get started locally, make sure you have Git and Go installed, then pull the
repository.

### Download Project

Pull repository with Git.

```bash
git pull https://github.com/pwilliams-ck/sniplate
```

## Build and Run Project

### Docker and Make

You can use `make` to build and run the docker environment, the migrating part
is a work in progress. You should be able to run `make migrate` to migrate up,
stick with the CLI for migrating down. You can check out the `infra/` folder for
more info.

### App and DB as Services

`cd` into the project root directory and execute the following to compile and
run with a single command. Further down there is a
[Build for Remote Server](#build-for-remote-server) section as well. The
database setup is a work in progress, and will use Make as well.

```bash
go run ./cmd/api -env=local-build
```

You should now be able to run `curl` commands against `localhost:4200`.

```bash
curl -i localhost:4200/v1/healthcheck
```

For server configuration info, try running `go run ./cmd/api -help`.

### Set Up TLS

There is an option to run the server with TLS enabled. If you want to create a
development TLS certificate and key, and have Go installed, they include a handy
tool to create self-signed certificates with the `crypto/tls` package. `cd` into
the project root, `mkdir tls`, and `cd` again.

```bash
cd sniplate
mkdir tls
cd tls
```

Next, find the `generate_cert.go` tool's path on your local machine, and run it
from the `tls` directory. Here is the path for MacOS, and probably Linux. The
CLI flags are there for convenience, copy/pasta away.

```bash
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

Then `cd` back into the root directory, build, and run the application.

```bash
cd ..
go run ./cmd/api -tls=true
```

If you are using a reverse proxy like _nginx_ or _HAProxy_ then this will
encrypt the traffic. between the proxy and the Sniplate app.

### Build Binary

First, make code changes.

Then run the following, where `./test-api` is the name of the executable you are
building.

```bash
GOOS=linux GOARCH=amd64 go build -o ./test-api ./cmd/api
```

## Conclusion

This is a work in progress and will be updated regularly.
