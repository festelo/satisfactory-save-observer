package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"

	"github.com/festelo/satisfactory-save-observer/internal/saves/adapters/handler"
	"github.com/festelo/satisfactory-save-observer/internal/saves/adapters/repository"
	"github.com/festelo/satisfactory-save-observer/internal/saves/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	directory := os.Getenv("DIRECTORY")
	host := os.Getenv("HOST")
	mapUrl := os.Getenv("MAP_URL")
	koDataPath := os.Getenv("KO_DATA_PATH")

	fmt.Println("DIRECTORY:", directory)
	fmt.Println("HOST:", host)
	fmt.Println("MAP_URL:", mapUrl)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", SavesRoutes(directory, host, mapUrl, path.Join(koDataPath, "list-saves.html")))
	http.ListenAndServe(":3000", r)
}

func SavesRoutes(directory string, host string, mapUrl string, listSavesPath string) chi.Router {
	r := chi.NewRouter()

	listSavesTemplate := template.Must(template.ParseFiles(listSavesPath))

	savesHandler := handler.NewSavesHandler(
		*domain.NewSavesService(
			repository.NewFilesSaveRepository(directory),
			repository.NewSimpleUrlResolverRepository(host, mapUrl),
		),
		*listSavesTemplate,
	)
	r.Get("/", savesHandler.ListSaves)
	r.Get("/{name}", savesHandler.GetSave)
	r.Options("/{name}", savesHandler.Cors)
	r.Get("/latest", savesHandler.GetSaveLatest)
	r.Options("/latest", savesHandler.Cors)
	return r
}
