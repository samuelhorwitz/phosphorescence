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
	scriptVersionRouter := func(r chi.Router) {
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
		r.Use(middleware.Disable) // TODO remove this when ready
		r.Use(middleware.Session)
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
			r.Route("/version", scriptVersionRouter)
			r.Route("/versions", scriptVersionRouter)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthorizePrivateScriptActions)
				r.Post("/duplicate", phosphor.DuplicateScript)
				r.Put("/", phosphor.UpdateScript)
				r.Put("/publish", phosphor.PublishScript)
				r.Delete("/", phosphor.DeleteScript)
			})
		})
	}
	r.Route("/script", scriptRouter)
	r.Route("/scripts", scriptRouter)
	scriptChainVersionRouter := func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Paginate)
			// r.Get("/", phosphor.GetScriptChainVersions)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthorizePrivateScriptActions)
				// r.Get("/draft", phosphor.GetPrivateScriptChainVersions)
				// r.Get("/drafts", phosphor.GetPrivateScriptChainVersions)
			})
		})
		r.Route("/{scriptChainVersionID}", func(r chi.Router) {
			r.Use(middleware.AuthorizeReadScriptVersion)
			// r.Get("/", phosphor.GetScriptChainVersion)
			// r.Post("/fork", phosphor.ForkScriptChainVersion)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthorizePrivateScriptActions)
				// r.Post("/duplicate", phosphor.DuplicateScriptChainVersion)
				// r.Delete("/", phosphor.DeleteScriptChainVersion)
			})
		})
	}
	scriptChainRouter := func(r chi.Router) {
		r.Use(middleware.Disable) // TODO remove this when ready
		r.Use(middleware.Session)
		r.Use(middleware.AuthenticatedSession)
		// r.Post("/", phosphor.CreateScriptChain)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Paginate)
			// r.Get("/", phosphor.ListPublicScriptChains)
			// r.Get("/my", phosphor.ListCurrentUserScriptChains)
		})
		r.Route("/{scriptChainID}", func(r chi.Router) {
			r.Use(middleware.AuthorizeReadScript)
			// r.Get("/", phosphor.GetScriptChain)
			// r.Post("/fork", phosphor.ForkScriptChain)
			r.Route("/version", scriptChainVersionRouter)
			r.Route("/versions", scriptChainVersionRouter)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthorizePrivateScriptActions)
				// r.Post("/duplicate", phosphor.DuplicateScriptChain)
				// r.Put("/", phosphor.UpdateScriptChain)
				// r.Delete("/", phosphor.DeleteScriptChain)
			})
		})
	}
	r.Route("/script-chain", scriptChainRouter)
	r.Route("/script-chains", scriptChainRouter)
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
		r.Use(middleware.AuthorizePremiumSpotifyUser)
		r.Use(middleware.SpotifyLimiter)
		r.Get("/", phosphor.ListSpotifyDevices)
		r.Put("/{deviceID}", phosphor.TransferPlayback)
	}
	r.Route("/device", deviceRouter)
	r.Route("/devices", deviceRouter)
	r.Route("/search", func(r chi.Router) {
		r.Use(middleware.Session)
		r.Get("/{query}", phosphor.Search)
		r.Get("/tag/{tag}", phosphor.SearchTag)
		r.Get("/recommendation/{query}", phosphor.RecommendedQuery)
	})
	r.Route("/official", func(r chi.Router) {
		r.Use(middleware.Session)
		r.Use(middleware.AuthorizeOfficialAccount)
		r.Get("/playlist/{playlistID}", phosphor.MakeOfficialPlaylist)
	})
	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "User-agent: *\nDisallow: /\n")
	})
	return r
}
