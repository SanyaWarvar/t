package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/gomail.v2"
)

// Структура для парсинга JSON запроса
type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// Конфигурация Gmail SMTP
const (
	SMTPHost = "smtp.gmail.com"
	SMTPPort = 587
)

var (
	gmailUser     string
	gmailPassword string
)

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS (для тестирования)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Обрабатываем preflight OPTIONS запрос
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Проверяем метод запроса
	if r.Method != "POST" {
		http.Error(w, `{"error": "Only POST method is supported"}`, http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON тело запроса
	var emailReq EmailRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&emailReq); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Валидируем обязательные поля
	if emailReq.To == "" || emailReq.Message == "" {
		http.Error(w, `{"error": "Fields 'to' and 'message' are required"}`, http.StatusBadRequest)
		return
	}

	// Если тема не указана, устанавливаем значение по умолчанию
	if emailReq.Subject == "" {
		emailReq.Subject = "Message from Go Server"
	}

	// Создаем и настраиваем email сообщение
	message := gomail.NewMessage()
	message.SetHeader("From", gmailUser)
	message.SetHeader("To", emailReq.To)
	message.SetHeader("Subject", emailReq.Subject)
	message.SetBody("text/plain", emailReq.Message)

	// Настраиваем диалер SMTP
	dialer := gomail.NewDialer(SMTPHost, SMTPPort, gmailUser, gmailPassword)

	// Отправляем email
	if err := dialer.DialAndSend(message); err != nil {
		errorMsg := fmt.Sprintf(`{"error": "Failed to send email: %s"}`, err.Error())
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	response := map[string]string{
		"status":  "success",
		"message": "Email sent successfully",
		"to":      emailReq.To,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	// Получаем учетные данные из переменных окружения
	gmailUser = "templeapi88@gmail.com"
	gmailPassword = "rhds btok ryjh vhdj"

	if gmailUser == "" || gmailPassword == "" {
		log.Fatal("GMAIL_USER and GMAIL_PASSWORD environment variables must be set")
	}

	// Регистрируем обработчик
	http.HandleFunc("/send-email", sendEmailHandler)

	// Запускаем сервер
	port := ":80"
	log.Printf("Starting server on port %s", port)
	log.Printf("Use POST /send-email to send emails")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
