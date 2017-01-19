package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

type Config struct {
	keyfile  string
	certfile string
	port     int
}

type TlsListener struct {
	sessions map[int]net.Conn
	conf     *Config
}

func (l *TlsListener) Listen() {
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")

	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	listener, err := tls.Listen("tcp", ":"+strconv.Itoa(l.conf.port), config)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		l.OnIncommingConnection(conn)
	}
}

func (l *TlsListener) OnIncommingConnection(conn net.Conn) {
	var nb_sessions = len(l.sessions)
	l.sessions[nb_sessions] = conn
	log.Println(fmt.Sprintf("Session %d opened", nb_sessions))
}

func (l *TlsListener) OnConnectionClosed(id int) {
	conn, ok := l.sessions[id]
	if !ok {
		log.Println("Unknown session id", id)
	} else {
		conn.Close()
		delete(l.sessions, id)
	}
	log.Println(fmt.Sprintf("Session %d closed", id))
}

func (l *TlsListener) Interact(sessionId int) error {
	conn, ok := l.sessions[sessionId]
	if !ok {
		return errors.New("Unknow session id")
	}
	fmt.Println("Interacting with session", sessionId)
	lineReader, _ := readline.NewEx(&readline.Config{
		Stdout:    os.Stdout,
		EOFPrompt: "exit",
		AutoComplete: readline.NewPrefixCompleter(
			readline.PcItem("background"),
			readline.PcItem("exit")),
	})
	for {
		go l.ReadStream(conn, sessionId)
		line, err := lineReader.Readline()
		if err != nil {
			log.Println(err)
		}
		line = strings.TrimSpace(line)
		if line == "background" {
			break
		}
		conn.Write([]byte(line + "\n"))
		if line == "exit" {
			l.OnConnectionClosed(sessionId)
			break
		}
	}
	return nil
}

func ListSessions(tlsListener *TlsListener) func(string) []string {
	return func(line string) []string {
		ids := make([]string, len(tlsListener.sessions))
		for k, _ := range tlsListener.sessions {
			ids[k] = strconv.Itoa(k)
		}
		return ids
	}
}

func (l *TlsListener) ReadStream(conn net.Conn, id int) {
	for {
		var buf [512]byte
		read, err := conn.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				l.OnConnectionClosed(id)
			}
			break
		}
		fmt.Print(string(buf[:read]))
	}
}

var (
	port     = flag.Int("port", 4040, "The listenning port")
	keyfile  = flag.String("keyfile", "", "The X509 keyfile")
	certfile = flag.String("certfile", "", "The X509 certificate")
)

func main() {
	flag.Parse()

	tlsListener := &TlsListener{
		conf: &Config{
			port:     *port,
			keyfile:  *keyfile,
			certfile: *certfile,
		},
		sessions: make(map[int]net.Conn, 0),
	}
	completer := readline.NewPrefixCompleter(
		readline.PcItem("set"),
		readline.PcItem("sessions",
			readline.PcItem("-i",
				readline.PcItemDynamic(ListSessions(tlsListener))),
			readline.PcItem("-k")),
		readline.PcItem("exit"),
	)

	l, err := readline.NewEx(&readline.Config{
		Prompt:            "hershell>>",
		HistoryFile:       "/tmp/hershell.hist",
		InterruptPrompt:   "^C",
		AutoComplete:      completer,
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	go tlsListener.Listen()

	defer l.Close()
	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "sessions"):
			args := strings.Split(line, " ")
			switch len(args) {
			case 1:
				fmt.Println("Session ID\tRemote host")
				for k, v := range tlsListener.sessions {
					fmt.Println(fmt.Sprintf("%d\t\t%s", k, v.RemoteAddr().String()))
				}
			case 3:
				sessionId, _ := strconv.Atoi(args[2])
				if err := tlsListener.Interact(sessionId); err != nil {
					log.Println(err)
				}
			default:
				fmt.Println("Usage: session [-i|-k SESSIONID] ")
			}
		case line == "exit":
			l.Close()
		default:
			if len(line) > 0 {
				fmt.Println("Command not found")
			}
			break
		}
	}
}
