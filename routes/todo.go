package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	db "crafted.api/config"
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
	UserId     int64              `json:"userId" bson:"userId"`
}

type InsertResponse struct {
	Id interface{} `json:"id"`
}

type TodosResponse struct {
	Todos []Todo `json:"todos"`
}

var mongoConn *mongo.Client = db.ConnectMongo()

func GetTodos(w http.ResponseWriter, r *http.Request) {
	statusQuery := r.URL.Query().Get("status")

	filter := false
	getCompleted := false

	if statusQuery == "0" {
		filter = true

	} else if statusQuery == "1" {
		filter = true
		getCompleted = true
	}

	todoCollection := db.GetCollection(mongoConn, "todos")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	findCondition := bson.M{}

	if filter {
		if getCompleted {
			findCondition = bson.M{"isComplete": true}
		} else {
			findCondition = bson.M{"isComplete": false}
		}
	}

	cur, err := todoCollection.Find(ctx, findCondition)

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

	todoId, err := primitive.ObjectIDFromHex(vars["id"])

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todoCollection := db.GetCollection(mongoConn, "todos")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	var todo Todo

	filter := bson.M{"_id": todoId}

	findErr := todoCollection.FindOne(ctx, filter).Decode(&todo)

	if findErr == mongo.ErrNoDocuments {
		log.Println("record does not exist")
		http.Error(w, "record does not exist", http.StatusNotFound)
		return
	} else if findErr != nil {
		log.Println(findErr)
		http.Error(w, findErr.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.Println("Unable to decode body")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	todoCollection := db.GetCollection(mongoConn, "todos")
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
	var todo Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.Println("Unable to decode body")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todoCollection := db.GetCollection(mongoConn, "todos")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	filter := bson.M{"_id": todo.Id}

	update := bson.M{
		"$set": bson.M{
			"title":      todo.Title,
			"isComplete": todo.IsComplete},
	}

	updateRes := todoCollection.FindOneAndUpdate(ctx, filter, update)

	if updateRes.Err() != nil {
		log.Println(updateRes.Err().Error())
		http.Error(w, updateRes.Err().Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("Document updated")
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// retrieve request params
	vars := mux.Vars(r)

	todoId, err := primitive.ObjectIDFromHex(vars["id"])

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todosCollection := db.GetCollection(mongoConn, "todos")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{
		"_id": todoId,
	}

	delRes := todosCollection.FindOneAndDelete(ctx, filter)

	if delRes.Err() != nil {
		log.Println("Unable to delete item")
		log.Println(delRes.Err().Error())
		http.Error(w, delRes.Err().Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("Item deleted")
}
