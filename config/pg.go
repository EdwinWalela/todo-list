package db

import (
	"context"
	"log"
	"os"

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

func Execute(smt string) {

}

func Query(smt string) {

}
