package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var db *pgxpool.Pool
var rdb *redis.Client

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
}

func main() {
	// 1. Check Environment Variables (Docker sets these)
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost" // Fallback for running locally without Docker
	}
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	// 2. Connect to Postgres (Notice we use dbHost and the password from docker-compose)
	dbUrl := fmt.Sprintf("postgres://skhan:password123@%s:5432/skhan", dbHost)
	var err error

	// Add a tiny delay to give Postgres time to finish booting inside Docker
	time.Sleep(2 * time.Second)

	db, err = pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// 3. Connect to Redis (Notice we use redisHost)
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v\n", err)
	}

	// 4. Setup Router & Start Server
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", handleShorten)
	mux.HandleFunc("GET /{code}", handleRedirect)

	fmt.Println("Server listening on :8080 with Redis Caching enabled!")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	code := generateShortCode(6)

	_, err := db.Exec(context.Background(),
		"INSERT INTO urls (short_code, original_url) VALUES ($1, $2)",
		code, req.URL)

	if err != nil {
		log.Printf("DB Error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	resp := ShortenResponse{ShortCode: code}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	ctx := context.Background()
	var originalURL string

	cachedURL, err := rdb.Get(ctx, code).Result()

	if err == redis.Nil {
		fmt.Println("Cache Miss! Reading from Postgres...")
		err := db.QueryRow(ctx, "SELECT original_url FROM urls WHERE short_code=$1", code).Scan(&originalURL)

		if err == pgx.ErrNoRows {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("DB Error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		rdb.Set(ctx, code, originalURL, 24*time.Hour)

	} else if err != nil {
		log.Printf("Redis Error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		fmt.Println("Cache Hit! Served directly from RAM.")
		originalURL = cachedURL
	}

	go func() {
		_, _ = db.Exec(context.Background(),
			"UPDATE urls SET click_count = click_count + 1 WHERE short_code=$1", code)
	}()

	http.Redirect(w, r, originalURL, http.StatusFound)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
