package monitor

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Broker handles SSE connections
type Broker struct {
	Clients    map[chan string]bool
	NewClients chan chan string
	Closing    chan chan string
	Mu         sync.Mutex
}

func NewBroker() *Broker {
	b := &Broker{
		Clients:    make(map[chan string]bool),
		NewClients: make(chan (chan string)),
		Closing:    make(chan (chan string)),
	}
	go b.listen()
	return b
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.NewClients:
			b.Clients[s] = true
			log.Println("Client added")
		case s := <-b.Closing:
			delete(b.Clients, s)
			log.Println("Client removed")
		}
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan string)

	b.NewClients <- messageChan

	defer func() {
		b.Closing <- messageChan
	}()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-messageChan:
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				log.Printf("Error writing to writer: %v\n", err)
			}
			flusher.Flush()
		}
	}
}
