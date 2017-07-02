package main

import (
	"os"

	"log"
	"net/http"

	"bufio"

	"fmt"
	"net"

	"time"

	"github.com/gorilla/websocket"
	"github.com/jessevdk/go-flags"
	"github.com/mkevac/hmm/internal/browser"
)

var opts struct {
	Http        string `long:"http" description:"Address on which to listen to" default:"localhost:0"`
	NoHeader    bool   `short:"n" long:"noheader" description:"Do not expect header in the first line"`
	Verbose     bool   `short:"v" long:"verbose" description:"Increase verbosity"`
	Last        string `short:"l" long:"last" description:"Which period to show" default:"24h"`
	LastSeconds float64
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(w http.ResponseWriter, r *http.Request, input chan string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for line := range input {
		if err := conn.WriteJSON(line); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	var err error

	_, err = flags.Parse(&opts)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	lastDuration, err := time.ParseDuration(opts.Last)
	if err != nil {
		log.Fatalf("Error while parsing duration %v", opts.Last)
	}

	opts.LastSeconds = lastDuration.Seconds()

	var r = bufio.NewReader(os.Stdin)

	output := make(chan string, 1024)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := indexHtmlTemplate.Execute(w, nil); err != nil {
			log.Print("Error while executing index.html template")
		}
	})
	http.HandleFunc("/main.js", func(w http.ResponseWriter, r *http.Request) {
		if err := mainJsTemplate.Execute(w, &opts); err != nil {
			log.Print("Error while executing main.js template")
		}
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, output)
	})

	ln, err := net.Listen("tcp", opts.Http)
	if err != nil {
		log.Fatalf("Error while listening on '%s': %s", opts.Http, err)
	}

	log.Printf("Opening browser")

	if !browser.Open("http://" + ln.Addr().String()) {
		fmt.Fprintf(os.Stderr, "Hmm is listening on http://%s\n", ln.Addr().String())
	}

	go func() {
		if err := http.Serve(ln, nil); err != nil {
			log.Fatal("Error while serving http: ", err)
		}
	}()

	s := bufio.NewScanner(r)
	for s.Scan() {
		text := s.Text()
		if opts.Verbose {
			fmt.Println(text)
		}
		output <- text
	}
	if err := s.Err(); err != nil {
		log.Fatalf("Error while scanning file: %s", err)
	}

	select {}
}
