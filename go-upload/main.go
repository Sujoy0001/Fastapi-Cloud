package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var numbers []int

// generate 25 random numbers (1–100)
func generateNumbers() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 25; i++ {
		numbers = append(numbers, rand.Intn(100)+1)
	}
	log.Println("Generated Numbers:", numbers)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	if q == "" {
		http.Error(w, "Missing search query ?q=number", http.StatusBadRequest)
		return
	}

	searchNum, err := strconv.Atoi(q)
	if err != nil {
		http.Error(w, "Query must be a number", http.StatusBadRequest)
		return
	}

	var matches []int
	for _, num := range numbers {
		if num == searchNum {
			matches = append(matches, num)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

func main() {
	generateNumbers()

	http.HandleFunc("/search", searchHandler)

	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
