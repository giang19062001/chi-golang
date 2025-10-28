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

func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	// r.Body là 'body request' nhưng dưới dạng stream
	log.Println("Body", r.Body)

	// tạo decoder JSON để parse dữ liệu trực tiếp từ stream
	decoder := json.NewDecoder(r.Body)
	log.Println("decoder", decoder)

	// tạo biến params kiểu struct để lưu dữ liệu đã parse được
	params := parameter{}
	log.Println("params", params)
	log.Println("&params", &params)

	// parse dữ liệu  và gán giá trị thực tế từ 'stream body request' vào vùng nhớ của 'params struct'
	err := decoder.Decode(&params)
	if err != nil {
		responseWithErr(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	log.Println("params.Name", params.Name)
	log.Println("&params.Name", &params.Name)

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		responseWithErr(w, 400, fmt.Sprintf("Error creating feed: %v", err))
		return
	}
	respondWithJSON(w, 201, changeStructFeedToFeed(feed))

}

func (apiCfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		responseWithErr(w, 400, fmt.Sprintf("Error creating feeds: %v", err))
		return
	}
	respondWithJSON(w, 201, changeStructFeedsToFeeds(feeds))

}
