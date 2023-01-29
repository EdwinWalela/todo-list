package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"crafted.api/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	Ctx            *context.Context
	TodoCollection *mongo.Collection
}

func main() {
	godotenv.Load()
	PORT, _ := strconv.Atoi(os.Getenv("PORT"))

	URL := fmt.Sprintf("0.0.0.0:%d", PORT)

	// Setup Mux router

	r := mux.NewRouter()

	r.HandleFunc("/todos", routes.CreateTodo).Methods("POST")
	r.HandleFunc("/todos", routes.GetTodos).Methods("GET")
	r.HandleFunc("/todos", routes.GetTodos).Methods("GET").Queries("status", "{status}")
	r.HandleFunc("/todos", routes.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", routes.DeleteTodo).Methods("DELETE")
	r.HandleFunc("/auth/register", routes.HandleRegister).Methods("POST")
	r.HandleFunc("/auth/login", routes.HandleLogin).Methods("POST")

	srv := &http.Server{
		Addr:         URL,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Printf("Listening for requests on port :%d\n", PORT)
	log.Fatal(srv.ListenAndServe())
}
