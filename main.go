package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/giang19062001/chi-golang/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// load .env file
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		defer fmt.Println("---")
		log.Fatal("Port not found") // ? in log ( có thời gian, tự thoát chương trình, ko chạy cả defer )
		// panic("Port not found")  // ? in log ( tự thoát chương trình, chạy cả defer )
	}

	// kiểm tra kết nối database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found")
	}

	conn, errDb := sql.Open("postgres", dbURL) // postgres, mysql, sqlite3...
	if errDb != nil {
		log.Fatal("Can't connect to database:", errDb)
	}

	// cấu hình api
	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	// cấu hình router
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReadiness)
	v1Router.Get("/err", handleErr)
	v1Router.Post("/users", apiCfg.handleCreateUser)
	router.Mount("/v1", v1Router)

	// cấu hình port và router cho server
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting on port %v", port)
	// kiểm tra nếu Server có lỗi
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
