package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"exploding-kittens-game/backend/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type Player struct {
	Name        string `json:"name"`
	Score       int    `json:"score"`
	GamesPlayed int    `json:"gamesPlayed"`
}

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// List of WebSocket clients
var clients = make(map[*websocket.Conn]bool)
var clientsMutex = sync.Mutex{}

// Entry point for server
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	utils.InitRedis()

	r := mux.NewRouter()
	r.HandleFunc("/updateScore", updateScoreHandler).Methods("POST")
	r.HandleFunc("/updateGames", updateGamesHandler).Methods("POST")
	r.HandleFunc("/leaderboard", leaderboardHandler).Methods("GET")
	r.HandleFunc("/ws-leaderboard", websocketHandler).Methods("GET")

	serverPort := os.Getenv("SERVER_PORT")
	log.Println("Server running on port", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, r))
}

// WebSocket handler to manage leaderboard updates
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	// Keep connection alive and listen for incoming messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			break
		}
	}
}

// Broadcast updated leaderboard to WebSocket clients
func broadcastLeaderboard() {
	leaderboard := fetchLeaderboard()

	data, err := json.Marshal(leaderboard)
	if err != nil {
		log.Println("Error marshalling leaderboard:", err)
		return
	}

	clientsMutex.Lock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Error sending message to client:", err)
			client.Close()
			delete(clients, client)
		}
	}
	clientsMutex.Unlock()
}

// Update score and broadcast leaderboard
func updateScoreHandler(w http.ResponseWriter, r *http.Request) {
	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	existingScore, err := utils.RedisClient.Get(utils.Ctx, "score:"+player.Name).Result()
	if err != nil {
		existingScore = "0"
	}

	newScore, _ := strconv.Atoi(existingScore)
	newScore += player.Score

	err = utils.RedisClient.Set(utils.Ctx, "score:"+player.Name, newScore, 0).Err()
	if err != nil {
		http.Error(w, "Failed to update score", http.StatusInternalServerError)
		return
	}

	// Broadcast updated leaderboard
	broadcastLeaderboard()

	w.WriteHeader(http.StatusOK)
}

// Update games played and broadcast leaderboard
func updateGamesHandler(w http.ResponseWriter, r *http.Request) {
	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	existingGames, err := utils.RedisClient.Get(utils.Ctx, "games:"+player.Name).Result()
	if err != nil {
		existingGames = "0"
	}

	newGames, _ := strconv.Atoi(existingGames)
	newGames += player.GamesPlayed

	err = utils.RedisClient.Set(utils.Ctx, "games:"+player.Name, newGames, 0).Err()
	if err != nil {
		http.Error(w, "Failed to update games played", http.StatusInternalServerError)
		return
	}

	// Broadcast updated leaderboard
	broadcastLeaderboard()

	w.WriteHeader(http.StatusOK)
}

// Fetch leaderboard from Redis
func fetchLeaderboard() []Player {
	keys, err := utils.RedisClient.Keys(utils.Ctx, "score:*").Result()
	if err != nil {
		log.Println("Error fetching leaderboard keys:", err)
		return nil
	}

	var leaderboard []Player
	for _, key := range keys {
		score, _ := utils.RedisClient.Get(utils.Ctx, key).Result()
		games, _ := utils.RedisClient.Get(utils.Ctx, "games:"+key[len("score:"):]).Result()
		player := Player{
			Name:        key[len("score:"):],
			Score:       atoi(score),
			GamesPlayed: atoi(games),
		}
		leaderboard = append(leaderboard, player)
	}
	return leaderboard
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	leaderboard := fetchLeaderboard()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
