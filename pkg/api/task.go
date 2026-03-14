package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"final-project/pkg/db"
	"final-project/pkg/logic"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch tasks"}); err != nil { // ✅ Ошибка обрабатывается
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(TasksResp{Tasks: tasks}); err != nil {
		log.Printf("[ERROR] Failed to encode response: %v", err)
		return
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "id is required"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "task not found"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("[ERROR] Failed to encode response: %v", err)
		return
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON format"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	if task.Title == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "title is required"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	if err := checkDate(&task); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "failed to save task"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(map[string]string{"id": strconv.FormatInt(id, 10)}); err != nil {
		log.Printf("[ERROR] Failed to encode response: %v", err)
		return
	}
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON format"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	if task.ID == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "id is required"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	if task.Title == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "title is required"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	if err := checkDate(&task); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	err := db.UpdateTask(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "task not found"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		log.Printf("[ERROR] Failed to encode response: %v", err)
		return
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "id is required"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "task not found"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		log.Printf("[ERROR] Failed to encode response: %v", err)
		return
	}
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "id is required"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)

		if err := json.NewEncoder(w).Encode(map[string]string{"error": "task not found"}); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
			return
		}
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)

			if err := json.NewEncoder(w).Encode(map[string]string{"error": "task not found"}); err != nil {
				log.Printf("[ERROR] Failed to encode response: %v", err)
				return
			}
			return
		}
	} else {
		next, err := logic.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)

			if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
				log.Printf("[ERROR] Failed to encode response: %v", err)
				return
			}
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)

			if err := json.NewEncoder(w).Encode(map[string]string{"error": "task not found"}); err != nil {
				log.Printf("[ERROR] Failed to encode response: %v", err)
				return
			}
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		log.Printf("[ERROR] Failed to encode response: %v", err)
		return
	}
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(logic.DateFormat)
		return nil
	}

	parsedDate, err := time.Parse(logic.DateFormat, task.Date)
	if err != nil {
		return err
	}

	if afterNow(now, parsedDate) {
		if task.Repeat == "" {
			task.Date = now.Format(logic.DateFormat)
		} else {
			next, err := logic.NextDate(now, task.Date, task.Repeat)

			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	return nil
}

func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return dateOnly.After(nowOnly)
}
