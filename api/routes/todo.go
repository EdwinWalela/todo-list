package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	db "crafted.api/config"
	"crafted.api/models"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InsertResponse struct {
	Id interface{} `json:"id"`
}

type TodosResponse struct {
	Todos []models.Todo `json:"todos"`
}

var mongoConn *mongo.Client = db.ConnectMongo()

func GetTodos(w http.ResponseWriter, r *http.Request) {
	headerToken := r.Header.Get("Auth")

	if len(headerToken) == 0 {
		log.Println("Token missing")
		http.Error(w, "Token missing", http.StatusBadRequest)
		return
	}

	// Extract & decode token
	token, err := jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		log.Println("Unable to decode token")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		log.Println("Token invalid")
		http.Error(w, "Token invalid", http.StatusBadRequest)
		return
	}

	statusQuery := r.URL.Query().Get("status")

	findCondition := bson.M{}

	filter := false
	getCompleted := false

	isAdmin := int(claims["is_admin"].(float64)) == 1
	userId := claims["id"]

	if statusQuery == "0" {
		filter = true

	} else if statusQuery == "1" {
		filter = true
		getCompleted = true
	}

	todoCollection := db.GetCollection(mongoConn, "todos")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	log.Println(userId)
	if filter {
		if getCompleted { // get completed items
			if !isAdmin {
				findCondition = bson.M{
					"isComplete": true,
					"user_id":    userId,
				}
			} else {
				findCondition = bson.M{
					"isComplete": true,
				}
			}
		} else { // get inot completed items
			if !isAdmin {
				findCondition = bson.M{
					"isComplete": false,
					"user_id":    userId,
				}
			} else {
				findCondition = bson.M{"isComplete": false}
			}

		}
	} else {
		if !isAdmin {
			findCondition = bson.M{"user_id": userId}
		}
	}

	cur, err := todoCollection.Find(ctx, findCondition)

	defer cur.Close(ctx)

	if err != nil {
		log.Println("unable to find collection")
		log.Println(err)
	}

	var todos = []models.Todo{}

	for cur.Next(ctx) {
		var todo models.Todo

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

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	headerToken := r.Header.Get("Auth")

	if len(headerToken) == 0 {
		log.Println("Token missing")
		http.Error(w, "Token missing", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.Println("Unable to decode body")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract & decode token
	token, err := jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		log.Println("Unable to decode token")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		log.Println("Token invalid")
		http.Error(w, "Token invalid", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	todoCollection := db.GetCollection(mongoConn, "todos")

	res, err := todoCollection.InsertOne(ctx, bson.D{
		{"title", todo.Title},
		{"timestamp", todo.Timestamp},
		{"isComplete", false},
		{"user_id", claims["id"]},
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
	var todo models.Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.Println("Unable to decode body")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	headerToken := r.Header.Get("Auth")

	if len(headerToken) == 0 {
		log.Println("Token missing")
		http.Error(w, "Token missing", http.StatusBadRequest)
		return
	}

	// Extract & decode token
	token, err := jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		log.Println("Unable to decode token")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		log.Println("Token invalid")
		http.Error(w, "Token invalid", http.StatusForbidden)
		return
	}

	todoCollection := db.GetCollection(mongoConn, "todos")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	filter := bson.M{"_id": todo.Id, "user_id": int64(claims["id"].(float64))}

	update := bson.M{
		"$set": bson.M{
			"title":      todo.Title,
			"isComplete": todo.IsComplete},
	}

	updateRes := todoCollection.FindOneAndUpdate(ctx, filter, update)

	if updateRes.Err() != nil {
		if updateRes.Err() == mongo.ErrNoDocuments {
			log.Println("Not user's todo item")
			http.Error(w, "Not user's todo item", http.StatusForbidden)
			return
		}
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

	headerToken := r.Header.Get("Auth")

	if len(headerToken) == 0 {
		log.Println("Token missing")
		http.Error(w, "Token missing", http.StatusBadRequest)
		return
	}

	// Extract & decode token
	token, err := jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		log.Println("Unable to decode token")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		log.Println("Token invalid")
		http.Error(w, "Token invalid", http.StatusForbidden)
		return
	}

	isAdmin := int(claims["is_admin"].(float64)) == 1

	if !isAdmin {
		log.Println("Non admin attempt to delete item")
		http.Error(w, "Unauthorized", http.StatusForbidden)
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
