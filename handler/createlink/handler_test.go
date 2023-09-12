package createlink_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/edmarfelipe/aws-lambda/handler/createlink"
	"github.com/edmarfelipe/aws-lambda/storage"
	"github.com/stretchr/testify/assert"
)

func TestCreateLink(t *testing.T) {
	linkStorageMock := &storage.LinkStorageMock{}
	handler := createlink.NewHandler(linkStorageMock)

	t.Run("Should create a link and return it", func(t *testing.T) {
		linkStorageMock.On("Create", storage.Link{Hash: "7378mDnD", Title: "My Link", Original: "https://www.google.com"}).Return(nil)

		body := `{"title":"My Link","link":"https://www.google.com"}`
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, `{"link":"http://127.0.0.1/7378mDnD"}`, rr.Body.String())
	})

	t.Run("Should return an error if the link is empty", func(t *testing.T) {
		body := `{"title":"My Link","link":""}`
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"message":"link is empty"}`, rr.Body.String())
	})

	t.Run("Should return an error if the title is empty", func(t *testing.T) {
		body := `{"title":"","link":"https://www.google.com"}`
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"message":"title is empty"}`, rr.Body.String())
	})

	t.Run("Should return an error if method is not POST", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
		assert.JSONEq(t, `{"message":"method not allowed"}`, rr.Body.String())
	})

	t.Run("Should return an error if body is invalid", func(t *testing.T) {
		body := ``
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"message":"internal server error"}`, rr.Body.String())
	})

	t.Run("Should return an error if fail to create a link", func(t *testing.T) {
		linkStorageMock.On("Create", storage.Link{Hash: "FnoLIXda", Title: "My Link", Original: "https://www.mylink.com"}).Return(errors.New("unknown error"))

		body := `{"title":"My Link","link":"https://www.mylink.com"}`
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"message":"internal server error"}`, rr.Body.String())
	})
}
