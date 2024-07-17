package db

// User представляет модель данных пользователя
type User struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PassportNum string `json:"passport_number"`
	// Добавьте другие поля, если необходимо
}

// Task представляет модель данных трудозатраты
type Task struct {
	ID       int     `json:"id"`
	UserID   int     `json:"user_id"`
	TaskName string  `json:"task_name"`
	Hours    float64 `json:"hours"`
	// Добавьте другие поля, если необходимо
}
