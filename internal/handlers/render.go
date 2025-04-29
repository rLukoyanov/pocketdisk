package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"pocketdisk/internal/models"
)

type Handlers struct {
	TemplateFS embed.FS
}

func (h *Handlers) Dashboard(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}
	log.Println(data)
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

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{"Message": "Хуй login"}

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
