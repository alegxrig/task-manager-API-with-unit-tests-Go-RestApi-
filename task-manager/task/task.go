package task

import (
	"errors"
	"fmt"
	"sync"
)

// Task defines the structure of our core data model
type Task struct {
	ID    int
	Title string
	Done  bool
}

// Manager handles concurrent-safe task operations using a Mutex
type Manager struct {
	mu     sync.Mutex
	tasks  []Task
	nextID int
}

// NewManager initializes and returns a pointer to a Manager instance
func NewManager() *Manager {
	return &Manager{
		tasks:  []Task{},
		nextID: 1,
	}
}

// AddTask safely appends a new task to the list
func (m *Manager) AddTask(title string) (Task, error) {
	if title == "" {
		return Task{}, errors.New("task title cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock() // Ensures the lock releases even if the function panics

	newTask := Task{
		ID:    m.nextID,
		Title: title,
		Done:  false,
	}
	m.nextID++

	m.tasks = append(m.tasks, newTask)
	return newTask, nil
}

// CompleteTask marks a task as done based on its unique ID
func (m *Manager) CompleteTask(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, t := range m.tasks {
		if t.ID == id {
			m.tasks[i].Done = true
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// ListTasks returns a copy of the tasks slice to prevent external modification
func (m *Manager) ListTasks() []Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return a copy of the slice to protect internal data state
	copiedTasks := make([]Task, len(m.tasks))
	copy(copiedTasks, m.tasks)
	return copiedTasks
}
