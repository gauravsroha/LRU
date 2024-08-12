package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Key        string `json:"key"`
		Value      string `json:"value"`
		Expiration int    `json:"expiration"` // in seconds
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data.Key == "" || data.Value == "" {
		http.Error(w, "Key and Value are required", http.StatusBadRequest)
		return
	}

	cache.Set(data.Key, data.Value, time.Duration(data.Expiration)*time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
