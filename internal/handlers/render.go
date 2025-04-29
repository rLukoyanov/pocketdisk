package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"pocketdisk/internal/config"
	"pocketdisk/internal/models"
)

type RenderHandlers struct {
	TemplateFS embed.FS
	Cfg        *config.Config
}

func (h *RenderHandlers) DashboardPage(w http.ResponseWriter, r *http.Request) {
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

func (h *RenderHandlers) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{"Error": "Хуй login"}

	t, err := template.ParseFS(h.TemplateFS, "templates/login.html")
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
