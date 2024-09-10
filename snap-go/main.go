package main

import (
	"fmt"
	"net/http"
)

func snapHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "game.html")
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Snap the game")
}

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/game", snapHandler)
	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
