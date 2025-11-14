package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/martencassel/opsmsg/catalog"
	"github.com/martencassel/opsmsg/dispatcher"
)

var todos []Todo
var nextID = 1

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func routes(d dispatcher.Dispatcher, c catalog.Catalog) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createTodoHandler(w, r, d, c)
		}
	})
	return mux
}

func createTodoHandler(w http.ResponseWriter, r *http.Request, d dispatcher.Dispatcher, c catalog.Catalog) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		msg := c.New("TODO002", map[string]string{"error": err.Error()})
		d.Dispatch(r.Context(), msg)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	todo.ID = nextID
	nextID++
	todos = append(todos, todo)

	msg := c.New("TODO001", map[string]string{"id": strconv.Itoa(todo.ID)})
	d.Dispatch(r.Context(), msg)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}
