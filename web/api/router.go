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
	allowedOrigins := []string{cfg.phosphorOrigin}
	if !cfg.isProduction {
		allowedOrigins = append(allowedOrigins, "http://localhost:8000")
	}
	cors := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)
	r.Use(middleware.CSP(cfg.phosphorOrigin))
	r.Use(chimiddleware.Timeout(cfg.handlerTimeout))
	r.Use(chimiddleware.NoCache)
	r.Use(chimiddleware.RealIP)
	r.Route("/spotify", func(r chi.Router) {
		r.Route("/authorize", func(r chi.Router) {
			r.Get("/", spotify.Authorize)
			r.Get("/redirect", spotify.AuthorizeRedirect)
		})
		r.Route("/unauthenticated", func(r chi.Router) {
			r.With(middleware.Captcha("api/tracks", 0.5)).Get("/tracks/{region}", spotify.TracksUnauthenticated)
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
			r.Post("/playlist", phosphor.CreateAndFollowPlaylist)
		})
	}
	r.Route("/user", userRouter)
	r.Route("/users", userRouter)
	trackRouter := func(r chi.Router) {
		r.Route("/unauthenticated", func(r chi.Router) {
			r.With(middleware.Captcha("api/track", 0.5)).Get("/{region}/{trackIDs}", phosphor.GetTracksUnauthenticated)
			r.With(middleware.Captcha("api/track/preview", 0.5)).Get("/preview/{region}/{trackIDs}", phosphor.GetTrackPreviewsUnauthenticated)
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.Session)
			r.Use(middleware.SpotifyLimiter)
			r.Get("/{trackIDs}", phosphor.GetTracks)
			r.Get("/preview/{trackIDs}", phosphor.GetTrackPreviews)
		})
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
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.Session)
		r.Use(middleware.AuthenticatedSession)
		r.Use(middleware.AuthorizeAdminAccount)
		r.Get("/playlist/{playlistID}", phosphor.MakePlaylistOfficial)
	})
	r.Route("/playlist", func(r chi.Router) {
		r.Route("/unauthenticated", func(r chi.Router) {
			r.With(middleware.Captcha("api/playlist", 0.5)).Get("/{region}/{playlistID}", phosphor.GetPlaylistUnauthenticated)
			r.With(middleware.Captcha("api/playlist/create", 0.5)).Post("/", phosphor.CreatePrivatePlaylist)
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.Session)
			r.Use(middleware.SpotifyLimiter)
			r.Get("/{playlistID}", phosphor.GetPlaylist)
			r.Post("/", phosphor.CreatePrivatePlaylist)
		})
	})
	r.Route("/album", func(r chi.Router) {
		r.Route("/unauthenticated", func(r chi.Router) {
			r.With(middleware.Captcha("api/album", 0.5)).Get("/{region}/{albumID}", phosphor.GetAlbumUnauthenticated)
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.Session)
			r.Use(middleware.SpotifyLimiter)
			r.Get("/{albumID}", phosphor.GetAlbum)
		})
	})
	r.Route("/player", func(r chi.Router) {
		r.With(middleware.Captcha("api/player/playlist", 0.5)).Get("/playlist/{playlistID}", phosphor.GetPlayerPlaylist)
	})
	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "User-agent: *\nDisallow: /\n")
	})
	return r
}
