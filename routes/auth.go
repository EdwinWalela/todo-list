package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	db "crafted.api/config"
	"crafted.api/models"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

var pgConn *pgx.Conn = db.ConnectPG()
var secret = []byte("1234")

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

	// Save user to db
	rows, err := db.CreateUser(pgConn, &user)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(fmt.Sprintf("%d row(s) affected", rows))
}

func genToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    user.Email,
		"is_admin": user.Is_admin,
		"nbf":      time.Date(2022, 2, 1, 12, 9, 9, 9, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("unable to decode request body")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user exists / password matches

	match, err := db.AuthenticateUser(pgConn, &user)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Email not registered", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !match {
		http.Error(w, "Incorrect combination", http.StatusForbidden)
		return
	}

	token, err := genToken(&user)

	if err != nil {
		log.Println("Unable to generate token")
		log.Println(err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(token)
}
