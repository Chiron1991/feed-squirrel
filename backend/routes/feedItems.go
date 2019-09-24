package routes

import (
	"github.com/Chiron1991/feed-squirrel/backend/controllers"
	"github.com/go-chi/chi"
)

func setupFeedItemRoutes(r chi.Router) {
	r.Get("/", controllers.FeedItemList)
	r.Route("/{feedItemID}", func(r chi.Router) {
		r.Use(controllers.FeedItemCtx)
		r.Get("/", controllers.FeedItemDetail)
	})
}
