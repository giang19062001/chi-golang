package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// kiểm tra biến port
	port := os.Getenv("PORT")
	if port == "" {
		defer fmt.Println("---")
		log.Fatal("Port not found") // ? in log ( có thời gian, tự thoát chương trình, ko chạy cả defer )
		// panic("Port not found")  // ? in log ( tự thoát chương trình, chạy cả defer )
	}

	// kiểm tra biến database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found")
	}

	// mở kết nối database
	conn, errDb := sql.Open("postgres", dbURL) // postgres, mysql, sqlite3...
	if errDb != nil {
		log.Fatal("Can't connect to database:", errDb)
	}

	// cấu hình api
	dbQueries := database.New(conn)
	apiCfg := apiConfig{
		DB: dbQueries,
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

	// chỉ định method cho router
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReadiness)
	v1Router.Get("/err", handleErr)
	// users
	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))
	// feeds
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.handleGetFeeds)
	// feed_follows
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	// cấu hình port và router cho server
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	// chạy ngầm gọi url lấy giá trị xml và log ra màn hình mỗi 1 phút từ 10 feed url có sẵn trong database
	const collectionConcurrency = 10
	const collectionInterval = time.Minute // 1 phút
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Server starting on port %v", port)
	// kiểm tra nếu Server có lỗi
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
