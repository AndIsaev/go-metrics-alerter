package metric

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"html/template"
	"net/http"
	"path/filepath"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	metrics := storage.MS.Metrics

	path := filepath.Join("internal/service/server/handler/metric", "html", "mainPage.html")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//выводим шаблон клиенту в браузер
	err = tmpl.Execute(w, metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
