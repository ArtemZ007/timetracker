package middleware

import (
	"log"
	"net/http"
	"time"
)

// InfoLog записывает информационное сообщение в стандартный вывод
func InfoLog(message string) {
	log.Printf("INFO: %s", message)
}

// ErrorLog записывает сообщение об ошибке в стандартный вывод ошибок
func ErrorLog(err error) {
	log.Printf("ERROR: %v", err)
}

// LoggingMiddleware является middleware для логирования запросов
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Запоминаем время начала обработки запроса

		// Обрабатываем запрос
		next.ServeHTTP(w, r)

		// Логируем информацию о запросе
		log.Printf("INFO: Метод: %s, URI: %s, Адрес: %s, Время обработки: %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}
