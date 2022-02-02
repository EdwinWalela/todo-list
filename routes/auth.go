package routes

type User struct {
	Id         int64  `json:"id"`
	Email      string `json:"email"`
	Password   int64  `json:"password"`
	IsComplete bool   `json:"isComplete"`
	IsAdmin    bool   `json:"isAdmin"`
}
