package Models

type User struct {
	UserID   uint   `json:"userID"`
	Username string `json:"username"`
	Password string `json:"password"`
}
