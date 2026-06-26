package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"task-manager/task"
)

// Глобальный менеджер задач
var tm = task.NewManager()

func main() {
	// Регистрируем маршруты (Endpoints)
	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/tasks/complete", handleCompleteTask)

	fmt.Println(" Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}

// handleTasks обрабатывает получение списка (GET) и создание (POST) задач
func handleTasks(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Возвращаем список всех задач
		tasks := tm.ListTasks()
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		// Структура для парсинга входящего JSON
		var body struct {
			Title string `json:"title"`
		}

		// Декодируем тело запроса
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		// Добавляем задачу через наш менеджер
		newTask, err := tm.AddTask(body.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Возвращаем созданную задачу со статусом 201 Created
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// handleCompleteTask обрабатывает завершение задачи через query-параметр (POST /tasks/complete?id=1)
func handleCompleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL (например, ?id=1)
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Параметр 'id' должен быть числом", http.StatusBadRequest)
		return
	}

	// Завершаем задачу
	if err := tm.CompleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Возвращаем успешный статус
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
