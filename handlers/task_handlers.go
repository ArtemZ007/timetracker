package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"timetracker/utils"
)

// TaskHandler представляет обработчик для задач
type TaskHandler struct {
	dbConn *utils.DB
}

// NewTaskHandler создает новый TaskHandler
func NewTaskHandler(dbConn *utils.DB) *TaskHandler {
	return &TaskHandler{dbConn: dbConn}
}

// GetTasks обрабатывает запрос на получение всех задач
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	// Получаем все задачи из базы данных
	tasks, err := h.dbConn.GetAllTasks()
	if err != nil {
		log.Printf("ERROR: Не удалось получить задачи: %v", err)
		http.Error(w, "Не удалось получить задачи", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Кодируем задачи в JSON и отправляем ответ
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Printf("ERROR: Не удалось закодировать задачи в JSON: %v", err)
		http.Error(w, "Не удалось закодировать задачи в JSON", http.StatusInternalServerError)
	}
}
