package main

import (
    "database/sql"
    "log"
    "os"
    "task-service/pkg/server"

    _ "github.com/lib/pq"
)

func main() {
    dbURL := os.Getenv("DATABASE_URL")
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatalf("failed to ping database: %v", err)
    }

    server.RunGRPCServer(db)
}
