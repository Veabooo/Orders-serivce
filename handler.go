package main

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// Главная — форма
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

// Отображение заказа
func orderHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	value, ok := cache.Load(id)
	if !ok {
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	order := value.(Order)
	tmpl.ExecuteTemplate(w, "order.html", order)
}
