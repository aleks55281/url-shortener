package handlers

import (
	"context"

	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

type Shortener interface {
	CreateShortUrl(ctx context.Context, origin string) (string, error)
	GetOriginUrl(ctx context.Context, shortUrl string) (string, error)
}
type Handler struct {
	shortener Shortener
	template  *template.Template
}

func NewHandler(short Shortener) *Handler {
	return &Handler{shortener: short}
}
func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/html/mainTemplate.html")
	if err != nil {
		slog.Error("template dont downoload", "error", err)
		http.Error(w, "template dont downoload", http.StatusBadRequest)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		slog.Error("template dont downoload")
		http.Error(w, "template dont downoload", http.StatusBadRequest)
		return
	}
}

func (h *Handler) Short(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("form dont parse", "error", err)
		http.Error(w, "form dont parse", http.StatusBadRequest)
	}

	origin := r.FormValue("url")

	shortUrl, err := h.shortener.CreateShortUrl(r.Context(), origin)
	if err != nil {
		slog.Error("couldn't shorted", "error", err)
		http.Error(w, "couldn't shorted", http.StatusBadRequest)
		return
	}
	tmpl, err := template.ParseFiles("static/html/PrintShortLink.html")
	if err != nil {
		slog.Error("template dont downoload", "error", err)
		http.Error(w, "template dont downoload", http.StatusBadRequest)
	}
	data := struct {
		ShortUrl string
	}{
		ShortUrl: shortUrl,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error("template dont downoload", "error", err)
		http.Error(w, "template dont downoload", http.StatusBadRequest)
	}
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 3 || path[2] == "" {
		slog.Error("invalid path", "path", r.URL.Path)
		http.Error(w, "invalid short path", http.StatusBadRequest)
		return
	}
	shortCode := path[2]

	originPath, err := h.shortener.GetOriginUrl(r.Context(), shortCode)
	if err != nil {
		slog.Error("invalid path", "path", r.URL.Path)
		http.Error(w, "invalid short path", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, originPath, http.StatusFound)

}
func (h *Handler) InitRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("/", http.HandlerFunc(h.MainPage))
	mux.Handle("/shorten", http.HandlerFunc(h.Short))
	mux.Handle("/sh/", http.HandlerFunc(h.RedirectHandler))
	return mux
}
