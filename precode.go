package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

const (
	contentTypeHeader     string = "Content-Type"
	jsonMIME              string = "application/json"
	tasksEndpointPattern  string = "/tasks"
	taskIdEndpointPattern string = "/task/{id}"
)

func getTasksHandler(w http.ResponseWriter, r *http.Request) {

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMIME)
	w.WriteHeader(http.StatusOK)

	w.Write(resp)
}

func postTasksHandler(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask

	w.Header().Set(contentTypeHeader, jsonMIME)
	w.WriteHeader(http.StatusCreated)

}

func getTaskIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMIME)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func deleteTaskIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, ok := tasks[id]; !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
	}
	delete(tasks, id)

	w.Header().Set(contentTypeHeader, jsonMIME)
	w.WriteHeader(http.StatusOK)

}

func main() {
	r := chi.NewRouter()

	r.Get(tasksEndpointPattern, getTasksHandler)
	r.Post(tasksEndpointPattern, postTasksHandler)
	r.Get(taskIdEndpointPattern, getTaskIdHandler)
	r.Delete(taskIdEndpointPattern, deleteTaskIdHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
