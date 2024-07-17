package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"timetracker/utils"
)

// UserHandler представляет обработчик для пользователей
type UserHandler struct {
	dbConn *utils.DB
}

// NewUserHandler создает новый UserHandler
func NewUserHandler(dbConn *utils.DB) *UserHandler {
	return &UserHandler{dbConn: dbConn}
}

// GetUsers обрабатывает запрос на получение всех пользователей
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Получаем всех пользователей из базы данных
	users, err := h.dbConn.GetAllUsers()
	if err != nil {
		log.Printf("ERROR: Не удалось получить пользователей: %v", err)
		http.Error(w, "Не удалось получить пользователей", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Кодируем пользователей в JSON и отправляем ответ
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("ERROR: Не удалось закодировать пользователей в JSON: %v", err)
		http.Error(w, "Не удалось закодировать пользователей в JSON", http.StatusInternalServerError)
	}
}
