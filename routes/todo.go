package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	mongoDriver "crafted.api/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	Title      string `json:"title"`
	Timestamp  int64  `json:"timestamp"`
	IsComplete bool   `json:"isComplete"`
}

type TodosResponse struct {
	Todos []Todo `json:"todos"`
}

var mongoConn *mongo.Client = mongoDriver.Connect()

func GetTodos(w http.ResponseWriter, r *http.Request) {

	todoCollection := mongoDriver.GetCollection(mongoConn, "todo")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	cur, err := todoCollection.Find(ctx, bson.D{})

	defer cur.Close(ctx)

	if err != nil {
		log.Println("unable to find collection")
		log.Println(err)
	}

	var todos = []Todo{}

	for cur.Next(ctx) {
		var todo Todo
		if err := cur.Decode(&todo); err != nil {
			log.Println("unable to decode item")
			log.Println(err)
		}

		todos = append(todos, todo)
	}

	res := &TodosResponse{
		Todos: todos,
	}

	json.NewEncoder(w).Encode(res)
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {

	// retrieve request params
	// vars := mux.Vars(r)

	// todoId := vars["id"]

	res := &Todo{
		Title:      "Ride bike",
		Timestamp:  time.Now().Unix(),
		IsComplete: false,
	}

	json.NewEncoder(w).Encode(res)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {

}
