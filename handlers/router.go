package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sse/events"
)

func InitRoutes(r *http.ServeMux) {
	handlerEvents := events.NewHandlerEvent()

	r.HandleFunc("/events", authMiddleware(handlerEvents.Handler))
	r.HandleFunc("/test1", authMiddleware(handlerTest1(handlerEvents)))
	r.HandleFunc("/test2", authMiddleware(handlerTest2(handlerEvents)))
	r.Handle("/", http.FileServer(http.Dir("./static")))
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cred := r.Header.Get("Authorization")
		if cred == "" {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		client := &http.Client{}
		req, _ := http.NewRequest("GET", "http://auth-api:8082/auth/verify", nil)
		req.Header.Set("Authorization", cred)

		resp, err := client.Do(req)
		if err != nil {
			log.Println(fmt.Printf("error connection to security server: %s", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		log.Println(fmt.Printf("security response: %s", resp.Body))

		next(w, r)
	}
}

func handlerTest1(n *events.HandlerEvent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data = map[string]any{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		n.Broadcast(events.EventMessage{
			EventName: "test1",
			Data:      data,
		})
	}
}

func handlerTest2(n *events.HandlerEvent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data = map[string]any{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		n.Broadcast(events.EventMessage{
			EventName: "test2",
			Data:      data,
		})
	}
}
