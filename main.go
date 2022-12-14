package main

import (
	"fmt"
	"log"
	"net/http"
	"sse/handlers"
)

func main() {
	r := http.NewServeMux()
	handlers.InitRoutes(r)
	err := http.ListenAndServe(":3500", r)

	if err != nil {
		log.Fatal(fmt.Printf("Server failed to start: %s", err))
	}

	log.Println("Server started")
}
