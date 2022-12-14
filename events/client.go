package events

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type client struct {
	ID          string
	sendMessage chan EventMessage
}

func newClient(id string) *client {
	return &client{
		ID:          id,
		sendMessage: make(chan EventMessage),
	}
}

func (c *client) Online(ctx context.Context, w io.Writer, flusher http.Flusher) {
	for {
		select {
		case m := <-c.sendMessage:
			data, err := json.Marshal(m.Data)
			if err != nil {
				log.Println(err)
			}
			const format = "event: %s\ndata: %s\n\n"
			_, err = fmt.Fprintf(w, format, m.EventName, string(data))
			if err != nil {
				log.Printf("Error writing to client: %s", err)
			}

			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}
