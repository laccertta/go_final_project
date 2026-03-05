package api

import (
 "net/http"
 "encoding/json"
 "time"

 "go_final_project/pkg/db"
)
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
 id := r.URL.Query().Get("id")
 if id == "" {
  writeJSON(w, map[string]string{"error": "Не указан идентификатор задачи"})
  return
 }
 err := db.DeleteTask(id)
 if err != nil {
  writeJSON(w, map[string]string{"error": err.Error()})
  return
 }
 writeJSON(w, map[string]string{})


 default:
  w.WriteHeader(http.StatusMethodNotAllowed)
  writeJSON(w, map[string]string{"error": "Метод не поддерживается"})
 }
}
// GET
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
 id := r.URL.Query().Get("id")
 task, err := db.GetTask(id)
 if err != nil {
  writeJSON(w, map[string]string{"error": err.Error()})
  return
 }

 writeJSON(w, task)
}

// PUT
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
 var task db.Task

 err := json.NewDecoder(r.Body).Decode(&task)
 if err != nil {
  writeJSON(w, map[string]string{"error": err.Error()})
  return
 }

 if task.ID == "" {
  writeJSON(w, map[string]string{"error": "Не указан идентификатор задачи"})
  return
 }

 if task.Title == "" {
  writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
  return
 }

 // проверка даты
 err = checkDate(&task)
 if err != nil {
  writeJSON(w, map[string]string{"error": err.Error()})
  return
 }

 // обновляем в БД
 err = db.UpdateTask(&task)
 if err != nil {
  writeJSON(w, map[string]string{"error": err.Error()})
  return
 }
 writeJSON(w, map[string]string{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
 id := r.URL.Query().Get("id")
 if id == "" {
  writeJSON(w, map[string]string{"error": "Не указан идентификатор задачи"})
  return
 }

 task, err := db.GetTask(id)
 if err != nil {
  writeJSON(w, map[string]string{"error": err.Error()})
  return
 }

 if task.Repeat == "" {
  // одноразовая задача,удаляем
  err := db.DeleteTask(id)
  if err != nil {
   writeJSON(w, map[string]string{"error": err.Error()})
   return
  }
  writeJSON(w, map[string]string{})
  return
 }

 // рассчитываем следующую дату
 next, err := NextDate(time.Now(), task.Date, task.Repeat)
 if err != nil {
  writeJSON(w, map[string]string{"error": err.Error()})
  return
 }

 err = db.UpdateDate(next, id)
 if err != nil {
  writeJSON(w, map[string]string{"error": err.Error()})
  return
 }

 writeJSON(w, map[string]string{})
}