package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/giang19062001/chi-golang/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		FeedID uuid.UUID `json:"feed_id"`
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

	log.Println("params after decode", params)

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.FeedID,
		UserID:    user.ID,
	})
	if err != nil {
		responseWithErr(w, 400, fmt.Sprintf("Error creating feed: %v", err))
		return
	}
	respondWithJSON(w, 201, changeStructFeedFollowToFeedFollow(feedFollow))

}

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		responseWithErr(w, 400, fmt.Sprintf("Error creating feed: %v", err))
		return
	}
	respondWithJSON(w, 201, changeStructFeedFollowsToFeedFollows(feedFollows))

}

func (apiCfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowId")
	log.Println("feedFollowIdStr", feedFollowIdStr)
	// convert dạng string sang dạng uuid như đã khai báo
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	log.Println("feedFollowId", feedFollowId)
	if err != nil {
		responseWithErr(w, 400, fmt.Sprintf("Cannot parse feed follow id : %v", err))
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		responseWithErr(w, 400, fmt.Sprintf("Cannot delete feed follow : %v", err))
	}
	respondWithJSON(w, 200, struct{}{})

}
