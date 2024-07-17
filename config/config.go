package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// AppConfig содержит настройки приложения
type AppConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	SecretKey  string
	ServerPort string
}

// LoadConfig загружает конфигурацию из файла .env
func LoadConfig() (*AppConfig, error) {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("ERROR: Не удалось загрузить файл .env: %v", err)
	}

	// Преобразуем значение порта базы данных из строки в целое число
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("ERROR: Неверное значение DB_PORT: %v", err)
	}

	// Создаем экземпляр AppConfig и заполняем его значениями из переменных окружения
	config := &AppConfig{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		SecretKey:  os.Getenv("SECRET_KEY"),
		ServerPort: os.Getenv("PORT"),
	}

	// Проверяем, что все необходимые переменные окружения заданы
	if config.DBHost == "" || config.DBUser == "" || config.DBPassword == "" || config.DBName == "" || config.SecretKey == "" || config.ServerPort == "" {
		return nil, fmt.Errorf("ERROR: Отсутствуют обязательные переменные окружения")
	}

	log.Println("INFO: Конфигурация успешно загружена из файла .env")
	return config, nil
}

// GetDatabaseURL возвращает строку подключения к базе данных
func (c *AppConfig) GetDatabaseURL() string {
	// Формируем строку подключения к базе данных PostgreSQL
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
