package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "crafted.api/config"
	"crafted.api/models"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

var pgConn *pgx.Conn = db.ConnectPG()

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("Unable to decode body")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate hash
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = string(hash)

	rows, err := db.CreateUser(pgConn, &user)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(fmt.Sprintf("%d row(s) affected", rows))
}
