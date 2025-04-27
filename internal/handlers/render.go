package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

type Handlers struct {
	TemplateFS embed.FS
}

func (h *Handlers) Render(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{"Message": "Хуй"}

	t, err := template.ParseFS(h.TemplateFS, "templates/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка рендера", http.StatusInternalServerError)
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка рендера", http.StatusInternalServerError)
	}
}
