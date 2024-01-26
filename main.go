package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	setupAPI()

	// generate using the `gencert.bash` script
	log.Fatal(http.ListenAndServeTLS(":8080", "server.cert", "server.key", nil))
}
func setupAPI() {

	manager := NewManager(context.Background())

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serverWS)
	http.HandleFunc("/login", manager.loginHandler)
}
