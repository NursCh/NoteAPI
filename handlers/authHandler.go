package handlers

import (
	"NotesAPI/Models"
	"NotesAPI/internal"
	"NotesAPI/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user Models.User

	db, err := internal.Connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	query := "SELECT ID FROM users WHERE username = $1 AND password = $2"
	var userID uint
	err = db.QueryRow(query, user.Username, user.Password).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateToken(userID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": token}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
