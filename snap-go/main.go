package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Card struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
}

var suits = []string{"Hearts", "Diamonds", "Clubs", "Spades"}
var values = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

// WebSocket upgrader to convert HTTP requests to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var players = make([]*websocket.Conn, 0, 2) // Store connections for both players
var deck []Card

// createDeck generates a deck of 52 cards
func createDeck() []Card {
	var deck []Card
	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, Card{Suit: suit, Value: value})
		}
	}
	return deck
}

// shuffleDeck shuffles the deck of cards
func shuffleDeck(deck []Card) []Card {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP request to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	// Add player to the players list
	players = append(players, conn)

	// If two players have connected, deal the deck
	if len(players) == 2 {
		dealDeckToPlayers(players)
	}

	// Keep the connection open and handle messages in a loop
	for {
		// Read message from the client
		_, _, err := conn.ReadMessage()
		if err != nil {
			// If an error occurs, remove the player and close the connection
			log.Println("Error reading message:", err)
			removePlayer(conn)
			conn.Close()
			break
		}
	}
}

// Remove player from the list
func removePlayer(conn *websocket.Conn) {
	for i, player := range players {
		if player == conn {
			players = append(players[:i], players[i+1:]...)
			break
		}
	}
}

// dealDeckToPlayers sends half the shuffled deck to each player
func dealDeckToPlayers(players []*websocket.Conn) {
	// Shuffle the deck and split it
	shuffledDeck := shuffleDeck(deck)
	halfDeckSize := len(shuffledDeck) / 2

	// Send half the deck to each player
	player1Deck := shuffledDeck[:halfDeckSize]
	player2Deck := shuffledDeck[halfDeckSize:]

	// Convert decks to JSON and send them to the players
	player1DeckJSON, _ := json.Marshal(player1Deck)
	player2DeckJSON, _ := json.Marshal(player2Deck)

	// Send deck to player 1
	err := players[0].WriteMessage(websocket.TextMessage, player1DeckJSON)
	if err != nil {
		log.Println("Error sending deck to player 1:", err)
		return
	}

	// Send deck to player 2
	err = players[1].WriteMessage(websocket.TextMessage, player2DeckJSON)
	if err != nil {
		log.Println("Error sending deck to player 2:", err)
		return
	}

	// Clear the players after the game starts
	players = players[:0]
}

func main() {
	// Create and shuffle the deck at startup
	deck = createDeck()

	// WebSocket endpoint for game
	http.HandleFunc("/ws", wsHandler)

	// Serve static frontend files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Serve the game.html file at the root
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/game.html")
	})

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}