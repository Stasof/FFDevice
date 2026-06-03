package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Структура для приёма данных из JSON
type UserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Обработчик для главной страницы — отдаёт HTML‑форму
func homeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	serial := query.Get("serial")
	check := query.Get("check")
	fmt.Println(serial, check)
	// Получаем отдельные значения
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "html/index.html")
		return
	}

	http.NotFound(w, r)
}

// Обработчик POST‑запросов с JSON
func submitHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса — POST
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем заголовок Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Ожидается Content-Type: application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Парсим JSON из тела запроса
	var userData UserData
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Простая валидация данных
	if userData.Name == "" || userData.Email == "" {
		http.Error(w, "Поля name и email обязательны", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "true")
}

func run() {
	// Регистрируем обработчики для путей
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api", submitHandler)

	// Запускаем сервер на порту 8080
	fmt.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
