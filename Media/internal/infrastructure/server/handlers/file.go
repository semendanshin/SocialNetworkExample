package handlers

import (
	"Media/internal/contracts/usecases"
	"Media/internal/infrastructure/server/utils/errorwrapper"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"mime/multipart"
	"net/http"
)

type FileHandler struct {
	fuc    usecases.FileUseCaseInterface
	logger *slog.Logger
}

func NewFileHandler(fuc usecases.FileUseCaseInterface, logger *slog.Logger) *FileHandler {
	return &FileHandler{
		fuc:    fuc,
		logger: logger,
	}
}

func (h *FileHandler) CreateFile(w http.ResponseWriter, r *http.Request) error {
	const op = "FileHandler.CreateFile"

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}

	files, ok := r.MultipartForm.File["file"]
	if !ok {
		return errors.New("file not found")
	}

	file, err := files[0].Open()
	if err != nil {
		h.logger.Error(op, slog.Any("error", err.Error()))
		return err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			h.logger.Error(op, slog.Any("error", err.Error()))
		}
	}(file)

	rawAuthorID := r.MultipartForm.Value["author_id"]
	if len(rawAuthorID) == 0 {
		return errors.New("author_id not found")
	}
	authorID, err := uuid.Parse(rawAuthorID[0])
	if err != nil {
		h.logger.Error(op, slog.Any("error", err.Error()))
		return err
	}
	dto := usecases.CreateFileDTO{
		Name:     files[0].Filename,
		AuthorID: authorID,
		Size:     files[0].Size,
		Content:  file,
	}

	id, err := h.fuc.CreateFile(r.Context(), dto)
	if err != nil {
		h.logger.Error(op, slog.Any("error", err.Error()))
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(`{"id": "` + id.String() + `"}`))

	if err != nil {
		h.logger.Error(op, slog.Any("error", err.Error()))
		return err
	}

	return nil
}

func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) error {
	const op = "FileHandler.GetFile"
	logger := h.logger.With("op", op)

	id := chi.URLParam(r, "id")

	fileID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("error while parsing file id", slog.Any("error", err.Error()))
		return err
	}

	file, err := h.fuc.GetFile(r.Context(), fileID)
	if err != nil {
		logger.Error("error while getting file", slog.Any("error", err.Error()))
		return err
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)

	_, err = w.Write(file.Content)

	return nil
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) error {
	const op = "FileHandler.DeleteFile"

	id := chi.URLParam(r, "id")

	fileID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error(op, slog.Any("error", err.Error()))
		return err
	}

	err = h.fuc.DeleteFile(r.Context(), fileID)
	if err != nil {
		h.logger.Error(op, slog.Any("error", err.Error()))
		return err
	}

	return nil
}

func (h *FileHandler) RegisterRoutes(mux *chi.Mux) {
	mux.Post("/files", errorwrapper.WrapWithError(h.CreateFile))
	mux.Get("/files/{id}", errorwrapper.WrapWithError(h.GetFile))
	mux.Delete("/files/{id}", errorwrapper.WrapWithError(h.DeleteFile))
}
