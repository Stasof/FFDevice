package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Структура для приёма данных из JSON
type BaseRESTRequest struct {
	Cmd  string `json:"cmd"`
	Args string `json:"args"`
}

type BaseRESTResponse struct {
	Status bool `json:"status"`
}

// Обработчик для главной страницы — отдаёт HTML‑форму
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем отдельные значения
	if r.URL.Path == "/" {
		query := r.URL.Query()
		ip := query.Get("ip")
		serial := query.Get("serial")
		check := query.Get("check")
		Init(ip, serial, check)
		fmt.Println(serial, check)
		http.ServeFile(w, r, "html/index.html")
		return
	}

	http.NotFound(w, r)
}

// Обработчик POST‑запросов с JSON
func submitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST")
	// Проверяем, что метод запроса — POST
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("POST1")
	// Проверяем заголовок Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Ожидается Content-Type: application/json", http.StatusUnsupportedMediaType)
		return
	}

	fmt.Println("POST2")
	// Парсим JSON из тела запроса
	var userData BaseRESTRequest
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("POST13")
	// Простая валидация данных
	if userData.Cmd == "" {
		http.Error(w, "Поля name и email обязательны", http.StatusBadRequest)
		return
	}

	response := BaseRESTResponse{Status: true}

	nresp := API(userData)
	if nresp != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(nresp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
