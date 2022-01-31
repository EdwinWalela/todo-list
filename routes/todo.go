package routes

import (
	"encoding/json"
	"net/http"
	"time"
)

type Todo struct {
	Title      string `json:"title"`
	Timestamp  int64  `json:"timestamp"`
	IsComplete bool   `json:"isComplete"`
}

type TodosResponse struct {
	Todos []Todo `json:"todos"`
}

func GetTodos(w http.ResponseWriter, r *http.Request) {

	todos := []Todo{
		{
			Title:      "Wash clothes",
			Timestamp:  time.Now().Unix(),
			IsComplete: false,
		},
		{
			Title:      "Read book",
			Timestamp:  time.Now().Unix(),
			IsComplete: true,
		},
		{
			Title:      "Ride bike",
			Timestamp:  time.Now().Unix(),
			IsComplete: false,
		},
	}

	res := &TodosResponse{
		Todos: todos,
	}

	json.NewEncoder(w).Encode(res)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {

}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {

}
