package routes

import (
	"github.com/Chiron1991/feed-squirrel/backend/controllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var r chi.Router

func init() {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.StripSlashes)
	r.Use(requestLoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/dev/load-fixtures", controllers.LoadDevFixtures) // todo: remove later
	r.Get("/dev/reset-db", controllers.ClearDB)              // todo: remove later
	r.Route("/feeds", setupFeedRoutes)
	r.Route("/feed-items", setupFeedItemRoutes)
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// wrap current ResponseWriter to gain access to the response
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		duration := time.Since(start)

		status := ww.Status()
		logFields := log.Fields{
			"requestID": middleware.GetReqID(r.Context()),
			"duration":  duration,
			"path":      r.URL,
			"method":    r.Method,
			"status":    ww.Status(),
			"remote":    r.RemoteAddr,
		}

		if status >= http.StatusInternalServerError {
			log.WithFields(logFields).Error("Error handling request")
		} else {
			log.WithFields(logFields).Info("Successfully handled request")
		}
	})
}

// GetRouter returns a router with all routes of the application
func GetRouter() *chi.Router {
	return &r
}
