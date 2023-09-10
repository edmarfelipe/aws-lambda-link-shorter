package redirectlink_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edmarfelipe/aws-lambda/handler/redirectlink"
	"github.com/edmarfelipe/aws-lambda/storage"
	"github.com/stretchr/testify/assert"
)

func TestRedirectLink(t *testing.T) {
	linkStorageMock := &storage.LinkStorageMock{}
	handler := redirectlink.NewHandler(linkStorageMock)

	t.Run("Should be a redirect link", func(t *testing.T) {
		linkStorageMock.On("GetLinkByHash", "12314").Return(storage.Link{
			Original: "https://www.google.com",
		}, nil)

		req, err := http.NewRequest(http.MethodGet, "/12314", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "https://www.google.com", rr.Header().Get("Location"))
	})

	t.Run("Should return 400 when link is empty", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{ "message": "link is empty" }`, rr.Body.String())
	})

	t.Run("Should return 404 when link not found", func(t *testing.T) {
		linkStorageMock.On("GetLinkByHash", "notfound").Return(storage.Link{}, storage.ErrLinkNotFound)

		req, err := http.NewRequest(http.MethodGet, "/notfound", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{ "message": "link not found" }`, rr.Body.String())
	})

	t.Run("Should return an error when internal server error", func(t *testing.T) {
		linkStorageMock.On("GetLinkByHash", "5555").Return(storage.Link{}, errors.New("unknown error"))

		req, err := http.NewRequest(http.MethodGet, "/5555", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{ "message": "internal server error" }`, rr.Body.String())
	})

	t.Run("Should return an error when method not allowed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
		assert.JSONEq(t, `{ "message": "method not allowed" }`, rr.Body.String())
	})
}
