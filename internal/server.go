package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	Port   string
	server *http.Server

	memory *Memory

	t  string
	cm chan string
}

func NewServer(port string, m *Memory, cm chan string) *Server {
	return &Server{Port: port, memory: m, cm: cm}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", s.callbackHandler)

	s.server = &http.Server{
		Addr:    ":" + s.Port,
		Handler: mux,
	}

	go func() {
		fmt.Printf("Server is running at http://localhost:%s\n", s.Port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan    // Block until an interrupt signal is received
	s.Shutdown() // Gracefully shut down the server
}

func (s *Server) Shutdown() {
	fmt.Println("\nShutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %s\n", err)
	} else {
		fmt.Println("Server stopped gracefully.")
	}
}

func (s *Server) callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "Missing code or state in query parameters", http.StatusBadRequest)
		return
	}
	s.cm <- code

	link := "http://localhost:8989/callback-to-server"
	params := url.Values{}
	params.Add("code", code)
	fullURL := link + "?" + params.Encode()

	resp, err := http.Get(fullURL)

	if err != nil {
		http.Error(w, "Error from the server", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error from the server", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	w.Write([]byte("Authorization successful! You can close this tab."))
}
