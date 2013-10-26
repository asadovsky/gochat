package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"code.google.com/p/go.net/websocket"
)

var port = flag.Int("port", 4000, "")

var listenAddr string

func init() {
	//listenAddr = fmt.Sprintf("0.0.0.0:%d", *port)
	listenAddr = fmt.Sprintf("localhost:%d", *port)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

type hub struct {
	clients     map[chan<- string]bool // set of active clients
	subscribe   chan chan<- string
	unsubscribe chan chan<- string
	broadcast   chan string
}

// TODO: Use buffer for broadcast channel?
func newHub() *hub {
	return &hub{
		clients:     make(map[chan<- string]bool),
		subscribe:   make(chan chan<- string),
		unsubscribe: make(chan chan<- string),
		broadcast:   make(chan string),
	}
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.subscribe:
			h.clients[c] = true
		case c := <-h.unsubscribe:
			delete(h.clients, c)
		case msg := <-h.broadcast:
			for send := range h.clients {
				send <- msg
			}
		}
	}
}

func wsHandler(h *hub, ws *websocket.Conn) {
	send := make(chan string)
	eof, done := make(chan bool), make(chan bool)
	h.subscribe <- send

	go func() {
		for {
			var msg string
			if err := websocket.Message.Receive(ws, &msg); err != nil {
				if err == io.EOF {
					eof <- true
					break
				}
				panicOnError(err)
			}
			h.broadcast <- msg
		}
	}()

	go func() {
	outer:
		for {
			select {
			case msg := <-send:
				panicOnError(websocket.Message.Send(ws, msg))
			case <-eof:
				break outer
			}
		}
		done <- true
	}()

	log.Printf("wsHandler WAIT, send=%v", send)
	<-done
	log.Printf("wsHandler EXIT, send=%v", send)

	// Clean up.
	h.unsubscribe <- send
	close(send)
	close(eof)
	close(done)
	ws.Close()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	if name == "" {
		name = "index.html"
	}

	ctype := mime.TypeByExtension(filepath.Ext(name))
	w.Header().Set("Content-Type", ctype)

	if name != "index.html" {
		b, err := ioutil.ReadFile(name)
		panicOnError(err)
		w.Write(b)
		return
	}

	tmpl := template.Must(template.ParseGlob("*.*"))
	panicOnError(tmpl.ExecuteTemplate(
		w, name, template.URL(fmt.Sprintf("ws://%s/ws", listenAddr))))
}

func main() {
	flag.Parse()
	h := newHub()
	go h.run()
	http.HandleFunc("/", rootHandler)
	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) { wsHandler(h, ws) }))
	log.Printf("Serving http://%s", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
