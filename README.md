# CocusProject

## Build and Install

For Build and Install, golang patterns are being used. See the link: https://go.dev/doc/

### Running

Golang allows you to run the project with the `run` command, which under the hood builds the project using your computer's tmpdir and then runs the binary.

to the Server:
```bash
$ go run cmd/server/main.go
```

to the Client:
```bash
$ go run cmd/client/main.go
```

The .env file is being used for the settings of each executable. You DO NOT need to change it if you run it locally.

### Redis Server 

If after running a Cocus Server, you keeping receive the message `Redis disconnected`. Check if redis server is running.

Running Redis localy:
```bash
$ redis-server
```

Access Redis localy:
```bash
$ redis-cli
```
