package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/scripts"
	"github.com/samuelhorwitz/phosphorescence/api/spotify"
	"net/http"
	"os"
	"path/filepath"
)

func initializeRoutes(cfg *config) http.Handler {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.phosphorOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	r.Use(cors.Handler)
	r.Use(middleware.CSP(cfg.phosphorOrigin))
	if !cfg.isProduction {
		r.Get("/tracks.json", serveTracks)
	}
	r.Route("/spotify", func(r chi.Router) {
		r.Get("/authorize", spotify.Authorize)
		r.Get("/tokens", spotify.Tokens)
		r.With(middleware.Authenticate).Get("/tracks", spotify.Tracks)
	})
	r.Route("/scripts", func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Get("/my", scripts.UserScripts)
		r.Post("/", scripts.SaveScript)
	})
	return r
}

func serveTracks(w http.ResponseWriter, r *http.Request) {
	ex, err := os.Executable()
	if err != nil {
		common.Fail(w, err, http.StatusInternalServerError)
		return
	}
	exPath := filepath.Dir(ex)
	tracksJSON := filepath.Join(exPath, "..", "static", "tracks.json")
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, tracksJSON)
}
