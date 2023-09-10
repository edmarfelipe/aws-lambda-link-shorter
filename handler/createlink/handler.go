package createlink

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/edmarfelipe/aws-lambda/storage"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrLinkEmpty        = errors.New("link is empty")
	ErrTitleEmpty       = errors.New("title is empty")
	ErrInternal         = errors.New("internal server error")
)

type Request struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type Response struct {
	Link string `json:"link"`
}

type ResponseError struct {
	Message string `json:"message"`
}

func (r *Request) Validate() error {
	if r.Title == "" {
		return ErrTitleEmpty
	}
	if r.Link == "" {
		return ErrLinkEmpty
	}
	return nil
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ResponseError{Message: ErrMethodNotAllowed.Error()})
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Message: ErrInternal.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{Message: err.Error()})
		return
	}

	link := storage.Link{
		Hash:     generateHashLink(req.Link),
		Title:    req.Title,
		Original: req.Link,
	}
	err = h.linkStorage.Create(r.Context(), link)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Message: ErrInternal.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Link: fmt.Sprintf("http://127.0.0.1/%s", link.Hash)})
}

func generateHashLink(originalLink string) string {
	hash := sha1.New()
	hash.Write([]byte(originalLink))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]
}
