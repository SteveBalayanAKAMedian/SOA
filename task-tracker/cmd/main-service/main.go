package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
	"tracker/pkg/user"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", "host=postgres port=5432 user=user password=password dbname=tasktracker sslmode=disable")
	if err != nil {
		log.Println("Failed to connect to database.")
		time.Sleep(5 * time.Second)
	}
	defer db.Close()

	userHandler := user.NewHandler(db)

	r := mux.NewRouter()
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/auth", userHandler.Authenticate).Methods("POST")

	r.Handle("/update", authMiddleware(db)(http.HandlerFunc(userHandler.Update))).Methods("PUT")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func authMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("session_token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := tokenCookie.Value
			session, err := user.GetSession(db, tokenStr)
			if err != nil || session.ExpiresAt.Before(time.Now()) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := user.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", session.UserID)
			ctx = context.WithValue(ctx, "username", claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
