package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// interface{}
// => ko có method nào -> một interface rỗng có thể chứa bất cứ type/struct nào trong Go (struct, int, string, map, slice, …)
// => nếu có method Chỉ chứa được các type/struct  có đủ method đó

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func responseWithErr(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		log.Println("Responding with 5xx error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}
