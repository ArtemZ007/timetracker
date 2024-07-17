package main

import (
	"log"
	"net/http"

	"timetracker/config"
	"timetracker/handlers"
	"timetracker/middleware"
	"timetracker/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		middleware.ErrorLog(err)
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	// Подключение к базе данных
	dbConn, err := utils.ConnectToDatabase(cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName)
	if err != nil {
		middleware.ErrorLog(err)
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			middleware.ErrorLog(err)
			log.Printf("Не удалось закрыть соединение с базой данных: %v", err)
		}
	}()

	// Инициализация роутера
	router := mux.NewRouter()

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(dbConn)
	taskHandler := handlers.NewTaskHandler(dbConn)

	// Определение маршрутов
	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")

	// Применение middleware
	router.Use(middleware.LoggingMiddleware)

	// Запуск сервера
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}
	middleware.InfoLog("Сервер запущен на порту " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		middleware.ErrorLog(err)
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
