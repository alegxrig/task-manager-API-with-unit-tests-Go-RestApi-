package task

import "testing"

// TestAddTask проверяет успешное добавление задачи и валидацию пустой строки
func TestAddTask(t *testing.T) {
	tm := NewManager()

	// Тест-кейс 1: Успешное добавление
	title := "Write unit tests"
	task, err := tm.AddTask(title)

	if err != nil {
		t.Fatalf("Ожидалась успешная генерация задачи, получена ошибка: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Неверный ID задачи. Ожидался: 1, получен: %d", task.ID)
	}

	if task.Title != title {
		t.Errorf("Неверное название. Ожидалось: %s, получено: %s", title, task.Title)
	}

	// Тест-кейс 2: Проверка валидации на пустую строку
	_, err = tm.AddTask("")
	if err == nil {
		t.Error("Ожидалась ошибка для пустого названия, но метод вернул nil")
	}
}

// TestCompleteTask проверяет смену статуса задачи и обработку несуществующих ID
func TestCompleteTask(t *testing.T) {
	tm := NewManager()
	task, _ := tm.AddTask("Test task")

	// Тест-кейс 1: Успешное завершение
	err := tm.CompleteTask(task.ID)
	if err != nil {
		t.Fatalf("Не удалось завершить задачу: %v", err)
	}

	tasks := tm.ListTasks()
	if !tasks[0].Done {
		t.Error("Статус задачи должен быть равен true (выполнено)")
	}

	// Тест-кейс 2: Попытка завершить несуществующий ID
	err = tm.CompleteTask(999)
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующего ID, но метод вернул nil")
	}
}
