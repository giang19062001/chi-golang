package main

import (
	"fmt"
	"net/http"

	"github.com/giang19062001/chi-golang/internal/auth"
	"github.com/giang19062001/chi-golang/internal/database"
)

// khai báo kiểu hàm authHandler với 3 tham số: Respons, Request, User
type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			responseWithErr(w, 403, fmt.Sprintf("Error getting API key: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			responseWithErr(w, 400, fmt.Sprintf("Error getting user: %v", err))
			return
		}

		// callback
		// truyền 3 tham số: Respons, Request, User vào authHandler
		handler(w, r, user)

	}
}
