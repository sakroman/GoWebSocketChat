package main

import (
	"fmt"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/chat.html")
}

func main() {
	manager := NewManager()

	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", manager.handleWebSocket)

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
