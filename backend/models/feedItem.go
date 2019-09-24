package models

import (
	"database/sql"
	"github.com/Chiron1991/feed-squirrel/backend/utils"
	"github.com/lib/pq"
)

// FeedItem represents an item within a Feed
type FeedItem struct {
	timestampedModel
	FeedID      int64       `db:"feed_id" json:"feedId"`
	Title       string      `db:"title" json:"title"`
	Description string      `db:"description" json:"description"`
	Content     string      `db:"content" json:"content"`
	Link        string      `db:"link" json:"link"`
	Updated     pq.NullTime `db:"updated" json:"updated"`
	Published   pq.NullTime `db:"published" json:"published"`
	//Author      sql.NullString `db:"author" json:"author"`
	GUID string `db:"guid" json:"guid"`
}

func GetAllFeedItems() *[]FeedItem {
	var feedItems []FeedItem
	query := `
		SELECT *
		FROM feed_items
		WHERE deleted_at IS NULL
		ORDER BY id ASC
	`
	if err := db.Select(&feedItems, query); err != nil {
		utils.LogSQLError(err, query)
	}
	return &feedItems
}

func GetFeedItemById(id int) *FeedItem {
	var feedItem FeedItem
	query := `
		SELECT *
		FROM feed_items
		WHERE deleted_at IS NULL AND id = $1
	`
	if err := db.Get(&feedItem, query, id); err != nil {
		utils.LogSQLError(err, query, id)
	}
	return &feedItem
}

func (feedItem FeedItem) Exists() bool {
	var result int
	query := `
		SELECT 1
		FROM feed_items
		WHERE guid = $1
		LIMIT 1
	`
	err := db.QueryRow(query, feedItem.GUID).Scan(&result)
	switch err {
	case nil:
		return true
	case sql.ErrNoRows:
		return false
	default:
		utils.LogSQLError(err, query, feedItem.GUID)
		return true
	}
}

// todo: this is shit, do better
func (feedItem FeedItem) Create() {
	query := `
		INSERT INTO feed_items (feed_id, title, description, link, guid)
		VALUES ($1, $2, $3, $4, $5)
	`
	tx := db.MustBegin()
	res := tx.MustExec(query, feedItem.FeedID, feedItem.Title, feedItem.Description, feedItem.Link, feedItem.GUID)
	if err := tx.Commit(); err != nil {
		utils.LogSQLError(err, query, feedItem.FeedID, feedItem.Title, feedItem.Description, feedItem.Link, feedItem.GUID)
	}
	feedItem.ID, _ = res.LastInsertId()
}
