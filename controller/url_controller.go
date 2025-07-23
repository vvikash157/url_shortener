package controller

import (
	"encoding/json"
	"net/http"
	"github.com/vvikash157/url_shortener/models"
	"github.com/vvikash157/url_shortener/services"

	"github.com/gorilla/mux"
)

type UrlController struct {
	services services.URLService
}

func NewUrlController(services services.URLService) *UrlController {
	return &UrlController{services: services}
}

func (s *UrlController) ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.LongUrl == "" {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	resp, err := s.services.UrlShortener(req)
	if err != nil {
		http.Error(w, "error while getting response from service UrlShortener", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *UrlController) RedirectOnURL(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	longUrl, err := s.services.ResolveUrl(code)
	if err != nil || longUrl == "" {
		http.NotFound(w, r)
	}

	http.Redirect(w, r, longUrl, http.StatusFound)
}
