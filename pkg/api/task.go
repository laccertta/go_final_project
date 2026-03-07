package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

// помощники основного хендлера
func getID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "task ID not specified")
		return "", false
	}
	return id, true
}

func writeJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// основной хендлер
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// вызываем функцию из addtask.go
		addTaskHandler(w, r)
	case http.MethodGet:
		// получаем задачу по id
		getTaskHandler(w, r)

	case http.MethodPut:
		// обновлем существующую задачу
		updateTaskHandler(w, r)

	case http.MethodDelete:
		id, ok := getID(w, r)
		if !ok {
			return
		}

		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, map[string]string{})

	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not supported")
	}
}

// GET
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(w, r)
	if !ok {
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, task)
}

// PUT
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.ID == "" {
		writeJSONError(w, http.StatusBadRequest, "task ID not specified")
		return
	}

	if task.Title == "" {
		writeJSONError(w, http.StatusBadRequest, "task title not specified")
		return
	}

	// проверка даты
	if err := checkDate(&task); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// обновляем в БД
	if err := db.UpdateTask(&task); err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, map[string]string{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(w, r)
	if !ok {
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	if task.Repeat == "" {
		// одноразовая задача,удаляем
		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	// рассчитываем следующую дату
	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, map[string]string{})
}
