package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// For listen port (mandatory)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env is required")
	}

	instanceID := os.Getenv("INSTANCE_ID")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Response handler
		if r.Method != "GET" {
			http.Error(w, "the requested http method is not allowed", http.StatusMethodNotAllowed)
		}
		// Response body
		text := "hello world!"
		if instanceID != "" {
			text = text + ". from " + instanceID
		}
		w.Write([]byte(text))
	})

	server := new(http.Server)
	server.Handler = mux
	server.Addr = ":" + port
	log.Println("web server is starting at", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}

}
