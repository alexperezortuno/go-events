package events

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type HandlerEvent struct {
	m       sync.Mutex
	clients map[string]*client
}

type EventMessage struct {
	EventName string
	Data      any
}

func NewHandlerEvent() *HandlerEvent {
	return &HandlerEvent{
		clients: make(map[string]*client),
	}
}

func (h *HandlerEvent) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	headers := r.Header
	log.Println(fmt.Printf("Headers: %s", headers))
	// ID with url params
	//id := r.URL.Query().Get("id")
	//if id == "" {
	//	fmt.Println("id is empty")
	//	id = "anonymous"
	//}

	// ID with header auth
	auth := r.Header.Get("Authorization")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	c := newClient(auth)
	h.register(c)
	fmt.Println("Client connected", auth)
	c.Online(r.Context(), w, flusher)
	fmt.Println("Client disconnected", auth)
	h.unregister(c)

	//flusher.Flush()
}

func (h *HandlerEvent) register(c *client) {
	h.m.Lock()
	defer h.m.Unlock()
	h.clients[c.ID] = c
}

func (h *HandlerEvent) unregister(c *client) {
	h.m.Lock()
	defer h.m.Unlock()
	delete(h.clients, c.ID)
}

func (h *HandlerEvent) Broadcast(m EventMessage) {
	h.m.Lock()
	defer h.m.Unlock()
	for _, c := range h.clients {
		c.sendMessage <- m
	}
}
