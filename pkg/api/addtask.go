package api

import (
 "encoding/json"
 "net/http"
 "time"

 "go_final_project/pkg/db"
)

// проверка даты задачи
func checkDate(task *db.Task) error {
 now := time.Now()
 todayStr := now.Format(dateFormat)

 if task.Date == "" {
  task.Date = todayStr
 }

 t, err := time.Parse(dateFormat, task.Date)
 if err != nil {
  return err
 }

 // если дата сегодня, оставляем её
 if sameDay(t, now) {
  task.Date = todayStr
  return nil
 }

 // если дата в прошлом
 if t.Before(now) {
  if task.Repeat == "" {
   task.Date = todayStr
  } else {
   next, err := NextDate(now, task.Date, task.Repeat)
   if err != nil {
    return err
   }
   task.Date = next
  }
 }

 return nil
}

func sameDay(a, b time.Time) bool {
 return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

// обработчик POST/api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
 var task db.Task

 // читаем JSON
 err := json.NewDecoder(r.Body).Decode(&task)
 if err != nil {
  writeJSON(w, map[string]string{"error": err.Error()})
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

 // добавление в БД
 id, err := db.AddTask(&task)
 if err != nil {
  writeJSON(w, map[string]string{"error": err.Error()})
  return
 }

 writeJSON(w, map[string]int64{"id": id})
}