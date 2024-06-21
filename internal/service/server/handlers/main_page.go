package handlers

import (
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"html/template"
	"net/http"
)

func MainPageHandler(mem *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		simpleTemplate := "" +
			"<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    " +
			"<meta charset=\"UTF-8\">\n    " +
			"<title>Example Template</title>\n</head>\n<body>\n<h1>Here are the keys and values from the map:" +
			"</h1>\n{{ range $key, $value := .}}\n<p>Key: {{ $key }}, Value: {{ $value }}" +
			"</p>\n{{ end }}\n</body>\n</html>"

		tmpl, err := template.New("test").Parse(simpleTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//выводим шаблон клиенту в браузер
		err = tmpl.Execute(w, mem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
