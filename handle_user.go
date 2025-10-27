package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/giang19062001/chi-golang/internal/database"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Name string `json:"name"`
	}

	// r.Body là 'body request' nhưng dưới dạng stream
	log.Println("Body", r.Body)
	decoder := json.NewDecoder(r.Body) // tạo decoder JSON để parse dữ liệu trực tiếp từ stream
	log.Println("decoder", decoder)

	params := parameter{}
	log.Println("params", params)
	log.Println("&params", &params)

	// &params là địa chỉ bộ nhớ của 'params struct'
	err := decoder.Decode(&params) // parse dữ liệu  và gán giá trị thực tế từ 'stream body request' vào vùng nhớ của 'params struct'
	if err != nil {
		responseWithErr(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	log.Println("params.Name", params.Name)
	log.Println("&params.Name", &params.Name)

	// r.Context() trả về context liên quan tới request này, có các thông tin:
	// Khi client ngắt kết nối, context này được cancel
	// Có thể dùng để timeout DB query hoặc hủy goroutine đang chạy
	// user là i bên trong hàm CreateUser()
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		responseWithErr(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}
	respondWithJSON(w, 200, user)

}
