package Models

import "time"

type Note struct {
	ID      uint      `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
	UserID  string    `json:"userID"`
}
