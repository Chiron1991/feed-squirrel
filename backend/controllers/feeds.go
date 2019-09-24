package controllers

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"

	"github.com/Chiron1991/feed-squirrel/backend/models"
)

// FeedCtx middleware annotates Feed to the Request's Context
func FeedCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: should be handled by router?!
		feedID, err := strconv.Atoi(chi.URLParam(r, "feedID"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		feed := models.GetFeedById(feedID)
		if feed.ID == 0 {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "feed", feed)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FeedList returns a JSON list of all Feeds
func FeedList(w http.ResponseWriter, r *http.Request) {
	feeds := models.GetAllFeeds()
	respondJSON(w, r, feeds)
}

// FeedDetail returns a JSON of a specific Feed
func FeedDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	feed, ok := ctx.Value("feed").(*models.Feed)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	respondJSON(w, r, feed)
}
