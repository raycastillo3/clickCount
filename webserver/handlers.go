package webserver

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"
)

func (a *APIServer) getClicksHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusBadRequest)
	}

	a.cache.mu.Lock()
	defer a.cache.mu.Unlock()

	js, err := json.Marshal(a.cache.values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (a *APIServer) clickCountHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, button string) {
	a.cache.mu.Lock()
	defer a.cache.mu.Unlock()

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/plain")
		count := int64(0)
		switch button {
		case "item":
			count = a.cache.values.Item
		case "addToCart":
			count = a.cache.values.AddToCart
		case "buy":
			count = a.cache.values.Buy
		}
		w.Write([]byte(strconv.FormatInt(count, 10)))
		return

	case http.MethodPut:
		switch button {
		case "item":
			a.cache.values.Item++
		case "addToCart":
			a.cache.values.AddToCart++
		case "buy":
			a.cache.values.Buy++
		}
		return

	default:
		http.Error(w, "only GET and PUT allowed", http.StatusBadRequest)
	}
}

func (a *APIServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join(a.webdir, "index.html")
	tmpl := template.Must(template.ParseFiles(indexPath))
	tmpl.Execute(w, struct {
		Name  string
		Items []string
	}{
		Name:  "Apache Datahub",
		Items: []string{"item", "addToCart", "buy"},
	})
}
