package db

import (
	"context"
	"log"
	"os"

	"crafted.api/models"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func ConnectPG() *pgx.Conn {
	godotenv.Load()
	PG_URI := os.Getenv("PG_DB_URL")

	conn, err := pgx.Connect(context.Background(), PG_URI)

	if err != nil {
		log.Fatal("Unable to connect to db", err.Error())
	}
	log.Println("Connected to PG db")
	return conn
}

func CreateUser(conn *pgx.Conn, user *models.User) (int64, error) {

	res, err := conn.Exec(context.Background(), "INSERT into users (email,password,is_admin) VALUES($1,$2,$3)", user.Email, user.Password, user.Is_admin)

	return res.RowsAffected(), err

}

func AuthenticateUser(conn *pgx.Conn, user *models.User) (bool, error) {
	var email string
	var password string

	err := conn.QueryRow(context.Background(), "SELECT email,password from users where email=$1", user.Email).Scan(&email, &password)

	if err != nil {
		return false, err
	}

	passwordMatch := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))

	if passwordMatch != nil {
		return false, nil
	}

	return true, nil

}
