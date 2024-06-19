package user

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int
	FirstName   string
	LastName    string
	Birthday    string
	Email       string
	PhoneNumber string
	Username    string
	Password    string
}

func (u *User) Create(db *sql.DB) error {
	err := db.QueryRow("SELECT id FROM users WHERE username=$1", u.Username).Scan()
	if err == nil {
		return errors.New("username already used")
	}

	err = db.QueryRow("INSERT INTO users (first_name, last_name, birthday, email, phone_number, username, password) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		u.FirstName, u.LastName, u.Birthday, u.Email, u.PhoneNumber, u.Username, u.Password).Scan(&u.ID)
	return err
}

func (u *User) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE users SET first_name=$1, last_name=$2, birthday=$3, email=$4, phone_number=$5 WHERE username=$6",
		u.FirstName, u.LastName, u.Birthday, u.Email, u.PhoneNumber, u.Username)
	return err
}

func Authenticate(db *sql.DB, username, password string) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
