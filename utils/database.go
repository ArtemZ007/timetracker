package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// DB представляет соединение с базой данных
type DB struct {
	*sql.DB
}

// Task представляет модель данных трудозатраты
type Task struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	TaskName    string    `json:"task_name"`
	Hours       float64   `json:"hours"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

// User представляет модель данных пользователя
type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PassportNum  string `json:"passport_number"`
	EnrichedInfo string `json:"enriched_info"`
}

// ConnectToDatabase подключается к базе данных PostgreSQL
func ConnectToDatabase(user, password, host, dbname string) (*DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("ERROR: Не удалось подключиться к базе данных: %v", err)
		return nil, err
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Printf("ERROR: Не удалось проверить соединение с базой данных: %v", err)
		return nil, err
	}

	log.Println("INFO: Успешное подключение к базе данных")
	return &DB{db}, nil
}

// CreateUser добавляет нового пользователя в базу данных
func (db *DB) CreateUser(user *User) error {
	query := `INSERT INTO users (first_name, last_name, passport_number, enriched_info) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.QueryRow(query, user.FirstName, user.LastName, user.PassportNum, user.EnrichedInfo).Scan(&user.ID)
	if err != nil {
		log.Printf("ERROR: Не удалось добавить пользователя: %v", err)
		return err
	}
	log.Printf("INFO: Пользователь добавлен с ID: %d", user.ID)
	return nil
}

// GetUserByID получает пользователя по его ID
func (db *DB) GetUserByID(userID int) (*User, error) {
	user := &User{}
	query := `SELECT id, first_name, last_name, passport_number, enriched_info FROM users WHERE id = $1`
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.PassportNum, &user.EnrichedInfo)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: Пользователь с ID %d не найден", userID)
			return nil, nil
		}
		log.Printf("ERROR: Не удалось получить пользователя: %v", err)
		return nil, err
	}
	return user, nil
}

// DeleteUser удаляет пользователя по его ID
func (db *DB) DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, userID)
	if err != nil {
		log.Printf("ERROR: Не удалось удалить пользователя: %v", err)
		return err
	}
	log.Printf("INFO: Пользователь с ID %d удален", userID)
	return nil
}

// UpdateUser обновляет данные пользователя
func (db *DB) UpdateUser(user *User) error {
	query := `UPDATE users SET first_name = $1, last_name = $2, passport_number = $3, enriched_info = $4 WHERE id = $5`
	_, err := db.Exec(query, user.FirstName, user.LastName, user.PassportNum, user.EnrichedInfo, user.ID)
	if err != nil {
		log.Printf("ERROR: Не удалось обновить данные пользователя: %v", err)
		return err
	}
	log.Printf("INFO: Данные пользователя с ID %d обновлены", user.ID)
	return nil
}

// GetTasksByUserID получает трудозатраты пользователя за период
func (db *DB) GetTasksByUserID(userID int, startDate, endDate time.Time) ([]Task, error) {
	query := `SELECT id, user_id, task_name, hours, start_time, end_time FROM tasks WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3 ORDER BY hours DESC`
	rows, err := db.Query(query, userID, startDate, endDate)
	if err != nil {
		log.Printf("ERROR: Не удало��ь получить трудозатраты: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.TaskName, &task.Hours, &task.StartTime, &task.EndTime); err != nil {
			log.Printf("ERROR: Не удалось прочитать данные задачи: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: Ошибка при чтении строк: %v", err)
		return nil, err
	}
	return tasks, nil
}

// StartTask начинает отсчет времени по задаче для пользователя
func (db *DB) StartTask(userID int, taskName string) (int, error) {
	query := `INSERT INTO tasks (user_id, task_name, start_time) VALUES ($1, $2, $3) RETURNING id`
	var taskID int
	err := db.QueryRow(query, userID, taskName, time.Now()).Scan(&taskID)
	if err != nil {
		log.Printf("ERROR: Не удалось начать задачу: %v", err)
		return 0, err
	}
	log.Printf("INFO: Задача начата с ID: %d", taskID)
	return taskID, nil
}

// EndTask заканчивает отсчет времени по задаче для пользователя
func (db *DB) EndTask(taskID int) error {
	query := `UPDATE tasks SET end_time = $1, hours = EXTRACT(EPOCH FROM (end_time - start_time)) / 3600 WHERE id = $2`
	_, err := db.Exec(query, time.Now(), taskID)
	if err != nil {
		log.Printf("ERROR: Не удалось закончить задачу: %v", err)
		return err
	}
	log.Printf("INFO: Задача с ID %d завершена", taskID)
	return nil
}

// GetAllTasks получает все задачи из базы данных
func (db *DB) GetAllTasks() ([]Task, error) {
	rows, err := db.Query("SELECT id, title, description FROM tasks")
	if err != nil {
		log.Printf("ERROR: Не удалось выполнить запрос к базе данных: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description); err != nil {
			log.Printf("ERROR: Не удалось сканировать строку: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ERROR: Ошибка при итерации по строкам: %v", err)
		return nil, err
	}

	return tasks, nil
}

// GetAllUsers получает всех пользователей из базы данных
func (db *DB) GetAllUsers() ([]User, error) {
	rows, err := db.Query("SELECT id, first_name, last_name, passport_number, enriched_info FROM users")
	if err != nil {
		log.Printf("ERROR: Не удалось выполнить запрос к базе данных: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.PassportNum, &user.EnrichedInfo); err != nil {
			log.Printf("ERROR: Не удалось сканировать строку: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ERROR: Ошибка при итерации по строкам: %v", err)
		return nil, err
	}

	return users, nil
}
