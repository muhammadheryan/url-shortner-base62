package transport

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/muhammadheryan/url-shortner-base62/application/url"
	"github.com/muhammadheryan/url-shortner-base62/constant"
	"github.com/muhammadheryan/url-shortner-base62/model"
	"github.com/muhammadheryan/url-shortner-base62/utils/errors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type RestHandler struct {
	URLApp url.URLApp
}

func NewTransport(URLApp url.URLApp) http.Handler {
	mux := mux.NewRouter()

	rh := &RestHandler{
		URLApp: URLApp,
	}

	// Swagger UI - setup sederhana
	mux.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API routes
	mux.HandleFunc("/url", rh.CreateURLShortner).Methods(http.MethodPost)
	mux.HandleFunc("/url/{shortURL}", rh.GetOriginalURL).Methods(http.MethodGet)

	return mux
}

// @Summary Create short URL
// @Description Create a new short URL from original URL
// @Accept json
// @Produce json
// @Param request body model.CreateURLShortnerRequest true "Create URL Request"
// @Success 200 {object} model.GetURLResponse
// @Failure 400 {object} errors.CustomError
// @Router /url [post]
func (s *RestHandler) CreateURLShortner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.CreateURLShortnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.SetCustomError(constant.ErrInvalidRequest))
		return
	}

	data, err := s.URLApp.CreateURLShortner(ctx, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeSuccess(w, data)
}

// @Summary Redirect to original URL
// @Description Redirect to original URL using short URL
// @Accept json
// @Produce json
// @Param shortURL path string true "Short URL"
// @Success 308 {string} string "Redirect to original URL"
// @Failure 404 {object} errors.CustomError
// @Router /url/{shortURL} [get]
func (s *RestHandler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get shortURL from URL parameter
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	if shortURL == "" {
		writeError(w, errors.SetCustomError(constant.ErrInvalidRequest))
		return
	}

	// Get original URL from database
	data, err := s.URLApp.GetURLByShortURL(ctx, shortURL)
	if err != nil {
		writeError(w, err)
		return
	}

	// Redirect to original URL with HTTP 308 (Permanent Redirect)
	w.Header().Set("Location", data.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
