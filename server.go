package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Структура для приёма данных из JSON
type BaseRESTRequest struct {
	Cmd     string `json:"cmd"`
	Args    string `json:"args"`
	Printer PrinterStruct
}

type PrinterStruct struct {
	IP     string `json:"ip"`
	Serial string `json:"serial"`
	Check  string `json:"check"`
}

type BaseRESTResponse struct {
	Status bool `json:"status"`
}

// Обработчик для главной страницы — отдаёт HTML‑форму
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем отдельные значения
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "html/index.html")
		return
	}
	http.ServeFile(w, r, "html"+r.URL.Path)
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
	var userData BaseRESTRequest
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
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

	port := ":8765"
	fmt.Println("Сервер запущен на http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
