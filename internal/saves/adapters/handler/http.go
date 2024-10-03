package handler

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	utils "github.com/festelo/satisfactory-save-observer/internal"
	"github.com/festelo/satisfactory-save-observer/internal/saves/domain"
)

type SavesHandler struct {
	service      domain.SavesService
	listTemplate template.Template
}

func NewSavesHandler(service domain.SavesService, listTemplate template.Template) *SavesHandler {
	return &SavesHandler{service, listTemplate}
}

func (h SavesHandler) ListSaves(w http.ResponseWriter, r *http.Request) {
	saves, err := h.service.ListSaves()

	if err != nil {
		slog.Error("Can't list saves", utils.ErrAttr(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	slog.Warn(saves[0].Link)
	h.listTemplate.Execute(w, saves)
}

func (h SavesHandler) GetSave(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	h.setCors(w)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%s.sav`, name))
	if err := h.service.CopySave(name, w); err != nil {
		slog.Error("Can't copy latest save", utils.ErrAttr(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h SavesHandler) GetSaveLatest(w http.ResponseWriter, r *http.Request) {
	h.setCors(w)
	w.Header().Set("Content-Disposition", `attachment; filename=latest.sav`)
	if err := h.service.CopyLatestSave(w); err != nil {
		slog.Error("Can't copy latest save", utils.ErrAttr(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h SavesHandler) Cors(w http.ResponseWriter, r *http.Request) {
	h.setCors(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h SavesHandler) setCors(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}
