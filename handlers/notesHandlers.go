package handlers

import (
	"NotesAPI/Models"
	"NotesAPI/internal"
	"NotesAPI/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func AddNotesHandler(w http.ResponseWriter, r *http.Request) {
	var note Models.Note
	userID := r.Context().Value("userID").(uint)
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := AddNote(note, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(note)
	if err != nil {
		return
	}

}

func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	notes, err := LoadNotes(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		return
	}
}

func LoadNotes(userID uint) ([]Models.Note, error) {
	db, err := internal.Connect()
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	query := "SELECT id, title, content, date, user_id FROM notes WHERE user_id = $1"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var notes []Models.Note
	for rows.Next() {
		var note Models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.Date, &note.UserID); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func AddNote(note Models.Note, userID uint) error {
	db, err := internal.Connect()
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	okTitle, _ := utils.CheckSpelling(note.Title)
	okContent, _ := utils.CheckSpelling(note.Content)
	if okTitle && okContent {
		query := `INSERT INTO notes (title, content, date, user_id) VALUES ($1, $2, $3, $4)`
		_, err = db.Exec(query, note.Title, note.Content, time.Now(), userID)
		if err != nil {
			return err
		}

		return nil
	} else {
		return errors.New("spelling errors detected")
	}
}
