package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	mongoDriver "crafted.api/config"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	Title      string             `json:"title" bson:"title"`
	Timestamp  int64              `json:"timestamp" bson:"timestamp"`
	IsComplete bool               `json:"isComplete" bson:"isComplete"`
}

type InsertResponse struct {
	Id interface{} `json:"id"`
}

type TodosResponse struct {
	Todos []Todo `json:"todos"`
}

var mongoConn *mongo.Client = mongoDriver.Connect()

func GetTodos(w http.ResponseWriter, r *http.Request) {

	todoCollection := mongoDriver.GetCollection(mongoConn, "todos")
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
	vars := mux.Vars(r)

	todoId := vars["id"]

	todoCollection := mongoDriver.GetCollection(mongoConn, "todos")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	var todo Todo

	err := todoCollection.FindOne(ctx, bson.D{{"id", todoId}}).Decode(&todo)

	if err == mongo.ErrNoDocuments {
		log.Println("record does not exist")
		http.Error(w, "record does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	todoCollection := mongoDriver.GetCollection(mongoConn, "todos")
	res, err := todoCollection.InsertOne(ctx, bson.D{
		{"title", todo.Title},
		{"timestamp", todo.Timestamp},
		{"isComplete", false},
	})

	if err != nil {
		log.Println("Unable to create new item")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	insertRes := &InsertResponse{
		Id: res.InsertedID,
	}

	json.NewEncoder(w).Encode(insertRes)

}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {

}
