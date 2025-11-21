package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// RequestPayload defines the expected JSON input from FastAPI
type RequestPayload struct {
	PartialNumber string `json:"partial_number"`
}

// ResponsePayload defines the JSON output sent back to FastAPI
type ResponsePayload struct {
	ChecksumDigit string `json:"checksum_digit"`
}

func main() {
	http.HandleFunc("/calculate-checksum", calculateChecksumHandler)

	fmt.Println("Checksum Engine (Go) running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Go server failed: %v\n", err)
	}
}

func calculateChecksumHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	partialNum := payload.PartialNumber
	if len(partialNum) != 15 {
		http.Error(w, `{"error": "Number must be exactly 15 digits"}`, http.StatusBadRequest)
		return
	}

	checksum := luhnCalculate(partialNum)

	response := ResponsePayload{ChecksumDigit: strconv.Itoa(checksum)}
	json.NewEncoder(w).Encode(response)
}

func luhnCalculate(partialNum string) int {
	sum := 0
	isDoubled := true

	for i := len(partialNum) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(partialNum[i]))

		if isDoubled {
			doubled := digit * 2
			if doubled > 9 {
				doubled = doubled%10 + doubled/10
			}
			sum += doubled
		} else {
			sum += digit
		}
		isDoubled = !isDoubled
	}

	checksum := (10 - (sum % 10)) % 10
	return checksum
}
