package controllers

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"

	"github.com/Chiron1991/feed-squirrel/backend/models"
)

// FeedItemCtx middleware annotates FeedItem to the Request's Context
func FeedItemCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: should be handled by router?!
		feedItemID, err := strconv.Atoi(chi.URLParam(r, "feedItemID"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		feedItem := models.GetFeedItemById(feedItemID)
		if feedItem.ID == 0 {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "feedItem", feedItem)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FeedItemList returns a JSON list of all FeedItems
func FeedItemList(w http.ResponseWriter, r *http.Request) {
	feedItems := models.GetAllFeedItems()
	respondJSON(w, r, feedItems)
}

// FeedItemDetail returns a JSON of a specific FeedItem
func FeedItemDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	feedItem, ok := ctx.Value("feedItem").(*models.FeedItem)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	respondJSON(w, r, feedItem)
}
