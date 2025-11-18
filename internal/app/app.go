package app

import (
	"github.com/aleks55281/url-shortener/internal/postgres"
	"github.com/aleks55281/url-shortener/internal/service"
	"github.com/aleks55281/url-shortener/internal/transport/http/handlers"
	datab "github.com/aleks55281/url-shortener/pkg/db"
	"log/slog"
	"net/http"
)

func RunServer() {
	postgr := datab.PostgrSql{
		Host:     "127.0.0.1",
		Port:     "5432",
		User:     "postgres",
		Dbname:   "urlshortener",
		Password: "goLANGninja",
		Sslmode:  "disable"}
	db, err := datab.ConPostgrSql(postgr)
	if err != nil {
		slog.Error("dont starts database", "error", err)
	}
	short := postgres.NewShortener(db)
	shortUrlService := service.NewShortenerUrl(short)
	handl := handlers.NewHandler(shortUrlService)
	mux := handl.InitRouter()
	slog.Info("server started")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("server dont started", "error", err)
	}

}
