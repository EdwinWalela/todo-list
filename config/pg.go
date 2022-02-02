package db

import (
	"context"
	"log"
	"os"

	"crafted.api/models"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
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

	res, err := conn.Exec(context.Background(), "INSERT into users (email,password,is_admin) VALUES($1,$2,$3)", user.Email, user.Password, user.IsAdmin)

	return res.RowsAffected(), err

}
