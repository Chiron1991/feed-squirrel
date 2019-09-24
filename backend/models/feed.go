package models

import (
	"github.com/Chiron1991/feed-squirrel/backend/utils"
	"github.com/lib/pq"
)

// Feed represents a remote source for FeedItems
type Feed struct {
	timestampedModel
	LastScraped pq.NullTime `db:"last_scraped" json:"lastScraped"`
	Title       string      `db:"title" json:"title"`
	Description string      `db:"description" json:"description"`
	Link        string      `db:"link" json:"link"`
	FeedLink    string      `db:"feed_link" json:"feedLink"`
	Updated     pq.NullTime `db:"updated" json:"updated"`
	Published   pq.NullTime `db:"published" json:"published"`
	Author      string      `db:"author" json:"author"`
	Language    string      `db:"language" json:"language"`
	Image       string      `db:"image" json:"image"`
	Copyright   string      `db:"copyright" json:"copyright"`
	Generator   string      `db:"generator" json:"generator"`
	Categories  string      `db:"categories" json:"categories"`
	// todo: cleanup stupid fields
}

func GetAllFeeds() *[]Feed {
	var feeds []Feed
	query := `
		SELECT *
		FROM feeds
		WHERE deleted_at IS NULL
		ORDER BY id ASC
	`
	if err := db.Select(&feeds, query); err != nil {
		utils.LogSQLError(err, query)
	}
	return &feeds
}

func GetFeedsDueForScraping() *[]Feed {
	var feeds []Feed
	query := `
		SELECT *
		FROM feeds
		WHERE deleted_at IS NULL
		AND (last_scraped IS NULL OR (NOW() at time zone 'utc' - last_scraped) > INTERVAL '15 minutes') 
	`
	if err := db.Select(&feeds, query); err != nil {
		utils.LogSQLError(err, query)
	}
	return &feeds
}

func GetFeedById(id int) *Feed {
	var feed Feed
	query := `
		SELECT *
		FROM feeds
		WHERE deleted_at IS NULL AND id = $1
	`
	if err := db.Get(&feed, query, id); err != nil {
		utils.LogSQLError(err, query, id)
	}
	return &feed
}

// todo: this is shit, do better
func (feed Feed) Create() {
	query := `
		INSERT INTO feeds (title, description, feed_link)
		VALUES ($1, $2, $3)
	`
	tx := db.MustBegin()
	res := tx.MustExec(query, feed.Title, feed.Description, feed.FeedLink)
	if err := tx.Commit(); err != nil {
		utils.LogSQLError(err, query, feed.Title, feed.Description, feed.FeedLink)
	}
	feed.ID, _ = res.LastInsertId()
}
