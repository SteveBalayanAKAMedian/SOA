package user

import (
    "database/sql"
    "time"
	"os"

    "github.com/golang-jwt/jwt/v4"
)

type Session struct {
    ID        int
    UserID    int
    Token     string
    ExpiresAt time.Time
}

func CreateSession(db *sql.DB, userID int, token string, expiresAt time.Time) error {
    _, err := db.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, token, expiresAt)
    return err
}

func GetSession(db *sql.DB, token string) (*Session, error) {
    session := &Session{}
    err := db.QueryRow("SELECT id, user_id, token, expires_at FROM sessions WHERE token = $1", token).Scan(
        &session.ID, &session.UserID, &session.Token, &session.ExpiresAt)
    if err != nil {
        return nil, err
    }
    return session, nil
}

func DeleteSession(db *sql.DB, token string) error {
    _, err := db.Exec("DELETE FROM sessions WHERE token = $1", token)
    return err
}

func ValidateToken(tokenStr string) (*jwt.RegisteredClaims, error) {
    signingKey := []byte(os.Getenv("JWT_SECRET_KEY"))

    token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return signingKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        return claims, nil
    } else {
        return nil, err
    }
}
