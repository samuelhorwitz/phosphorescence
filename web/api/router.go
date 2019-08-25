package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/samuelhorwitz/phosphorescence/api/handlers/phosphor"
	"github.com/samuelhorwitz/phosphorescence/api/handlers/spotify"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
)

func initializeRoutes(cfg *config) http.Handler {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.phosphorOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)
	r.Use(middleware.CSP(cfg.phosphorOrigin))
	r.Use(chimiddleware.Timeout(cfg.handlerTimeout))
	r.Use(chimiddleware.NoCache)
	r.Use(middleware.IPLimiter)
	r.Route("/spotify", func(r chi.Router) {
		r.Route("/authorize", func(r chi.Router) {
			r.Get("/", spotify.Authorize)
			r.Get("/redirect", spotify.AuthorizeRedirect)
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.Session)
			r.Get("/tracks", spotify.Tracks)
			r.Get("/token", spotify.Token)
		})
	})
	r.Route("/authenticate", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Session)
			r.Post("/", phosphor.Authenticate)
			r.Get("/logout", phosphor.Logout)
			r.Post("/logoutall", phosphor.LogoutEverywhere)
		})
		r.Get("/{magicLink}", phosphor.AuthenticateRedirect)
	})
	versionRouter := func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Paginate)
			r.Get("/", phosphor.GetScriptVersions)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthorizePrivateScriptActions)
				r.Get("/draft", phosphor.GetPrivateScriptVersions)
				r.Get("/drafts", phosphor.GetPrivateScriptVersions)
			})
		})
		r.Route("/{scriptVersionID}", func(r chi.Router) {
			r.Use(middleware.AuthorizeReadScriptVersion)
			r.Get("/", phosphor.GetScriptVersion)
			r.Post("/fork", phosphor.ForkScriptVersion)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthorizePrivateScriptActions)
				r.Post("/duplicate", phosphor.DuplicateScriptVersion)
				r.Delete("/", phosphor.DeleteScriptVersion)
			})
		})
	}
	scriptRouter := func(r chi.Router) {
		r.Use(middleware.Session)
		r.Get("/search", phosphor.Search)
		r.Get("/search-tag", phosphor.SearchTag)
		r.Get("/query-recommendation", phosphor.RecommendedQuery)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthenticatedSession)
			r.Post("/", phosphor.CreateScript)
			r.Group(func(r chi.Router) {
				r.Use(middleware.Paginate)
				r.Get("/", phosphor.ListPublicScripts)
				r.Get("/my", phosphor.ListCurrentUserScripts)
			})
			r.Route("/{scriptID}", func(r chi.Router) {
				r.Use(middleware.AuthorizeReadScript)
				r.Get("/", phosphor.GetScript)
				r.Post("/fork", phosphor.ForkScript)
				r.Route("/version", versionRouter)
				r.Route("/versions", versionRouter)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AuthorizePrivateScriptActions)
					r.Post("/duplicate", phosphor.DuplicateScript)
					r.Put("/", phosphor.UpdateScript)
					r.Put("/publish", phosphor.PublishScript)
					r.Delete("/", phosphor.DeleteScript)
				})
			})
		})
	}
	r.Route("/script", scriptRouter)
	r.Route("/scripts", scriptRouter)
	userRouter := func(r chi.Router) {
		r.Use(middleware.Session)
		r.Route("/me", func(r chi.Router) {
			r.Get("/", phosphor.GetCurrentUser)
			r.Get("/currently-playing", phosphor.GetCurrentlyPlaying)
			r.Post("/playlist", phosphor.CreatePlaylist)
		})
	}
	r.Route("/user", userRouter)
	r.Route("/users", userRouter)
	trackRouter := func(r chi.Router) {
		r.Use(middleware.Session)
		r.Use(middleware.SpotifyLimiter)
		r.Get("/{trackID}", phosphor.GetTrackData)
	}
	r.Route("/track", trackRouter)
	r.Route("/tracks", trackRouter)
	deviceRouter := func(r chi.Router) {
		r.Use(middleware.Session)
		r.Use(middleware.SpotifyLimiter)
		r.Get("/", phosphor.ListSpotifyDevices)
		r.Put("/{deviceID}", phosphor.TransferPlayback)
	}
	r.Route("/device", deviceRouter)
	r.Route("/devices", deviceRouter)
	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "User-agent: *\nDisallow: /\n")
	})
	return r
}
