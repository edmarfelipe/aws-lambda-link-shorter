package redirectlink

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/edmarfelipe/aws-lambda/storage"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrLinkNotFound     = errors.New("link not found")
	ErrLinkEmpty        = errors.New("link is empty")
	ErrInternal         = errors.New("internal server error")
)

type ResponseError struct {
	Message string `json:"message"`
}

type Handler struct {
	linkStorage storage.LinkStorage
}

func NewHandler(linkStorage storage.LinkStorage) *Handler {
	return &Handler{
		linkStorage: linkStorage,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ResponseError{Message: ErrMethodNotAllowed.Error()})
		return
	}

	if len(r.URL.Path) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{Message: ErrLinkEmpty.Error()})
		return
	}

	hash := r.URL.Path[1:]
	link, err := h.linkStorage.GetLinkByHash(r.Context(), hash)
	if err == storage.ErrLinkNotFound {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseError{Message: ErrLinkNotFound.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Message: ErrInternal.Error()})
		return
	}

	http.Redirect(w, r, link.Original, http.StatusTemporaryRedirect)
}
