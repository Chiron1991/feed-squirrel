package routes

import (
	"github.com/Chiron1991/feed-squirrel/backend/controllers"
	"github.com/go-chi/chi"
)

func setupFeedRoutes(r chi.Router) {
	r.Get("/", controllers.FeedList)
	r.Route("/{feedID}", func(r chi.Router) {
		r.Use(controllers.FeedCtx)
		r.Get("/", controllers.FeedDetail)
	})
}
