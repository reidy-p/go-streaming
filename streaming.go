package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)
func streamHandler(w http.ResponseWriter, r *http.Request) {
	// Create a ticker to send data at regular intervals
	ticker := time.NewTicker(1 * time.Second)
	bytesSent := 0

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			if bytesSent > 100 {
				return
			}

			buf := []byte(strings.Repeat("x", 4))
			n, err := w.Write(buf)
			bytesSent += n
			if err != nil {
				fmt.Println("Error writing to the client:", err)
				return
			}

			// Flush the response writer to ensure the data is sent immediately
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			fmt.Println("Client closed the connection")
			return
		}
	}
}

func startServer() {
	http.HandleFunc("/stream", streamHandler)

	fmt.Println("Streaming server started. Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func main() {
	go startServer()

	// default
	httpClient := http.Client{}

	// Timeout getting header
	//httpClient := http.Client{Timeout: 10 * time.Millisecond}

	// Timeout in getting response
	//httpClient := http.Client{Timeout: 2 * time.Second}

	resp, err := httpClient.Get("http://localhost:8080/stream")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 4)
	for {
		_, err := resp.Body.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		fmt.Println(buf)	
	}
}
