package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"example.com/notes-api/internal/core"
	"example.com/notes-api/internal/repo"
)

type Handler struct {
	Repo *repo.NoteRepoMem
}

// CreateNote godoc
// @Summary      Создать заметку
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        input  body     core.NoteCreate  true  "Данные новой заметки"
// @Success      201    {object} core.Note
// @Failure      400    {object} map[string]string
// @Failure      500    {object} map[string]string
// @Router       /notes [post]
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var input core.NoteCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	note := core.Note{
		Title:   input.Title,
		Content: input.Content,
	}

	id, err := h.Repo.Create(note)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	note.ID = id
	note.CreatedAt = time.Now().UTC()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// ListNotes godoc
// @Summary      Список заметок
// @Description  Возвращает список заметок
// @Tags         notes
// @Success      200    {array}  core.Note
// @Failure      500    {object} map[string]string
// @Router       /notes [get]
func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	notes := h.Repo.List()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// GetNote godoc
// @Summary      Получить заметку
// @Tags         notes
// @Param        id   path   int  true  "ID"
// @Success      200  {object}  core.Note
// @Failure      404  {object}  map[string]string
// @Router       /notes/{id} [get]
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	note, ok := h.Repo.Get(id)
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// PatchNote godoc
// @Summary      Обновить заметку (частично)
// @Tags         notes
// @Accept       json
// @Param        id     path   int        true  "ID"
// @Param        input  body   core.NoteUpdate true  "Поля для обновления"
// @Success      200    {object}  core.Note
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /notes/{id} [patch]
func (h *Handler) PatchNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var input core.NoteUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updated, ok := h.Repo.Update(id, input, time.Now().UTC())
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteNote godoc
// @Summary      Удалить заметку
// @Tags         notes
// @Param        id  path  int  true  "ID"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /notes/{id} [delete]
func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if !h.Repo.Delete(id) {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
