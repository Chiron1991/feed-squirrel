package main

import (
	"github.com/Chiron1991/feed-squirrel/backend/models"
	"github.com/Chiron1991/feed-squirrel/backend/utils"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
)

type UserAgentTransport struct {
	http.RoundTripper
}

func (c *UserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "FeedSquirrel / 1.0")
	return c.RoundTripper.RoundTrip(r)
}

func ScrapeDueFeeds() {
	// get due feeds from db
	dueFeeds := models.GetFeedsDueForScraping()

	jobs := make(chan models.Feed, 100)

	// declare workers
	for id := 1; id <= runtime.NumCPU(); id++ {
		go scrapeWorker(id, jobs)
	}

	// stuff feeds into worker
	for _, feed := range *dueFeeds {
		jobs <- feed
	}
	close(jobs)
}

func scrapeWorker(id int, jobs <-chan models.Feed) {
	for feed := range jobs {
		log.WithFields(log.Fields{
			"workerId": id,
			"feedID":   feed.ID,
			"feedLink": feed.FeedLink,
		}).Debug("Scraping feed")

		parser := gofeed.NewParser()
		// workaround for problems with User-Agent, see https://github.com/mmcdole/gofeed/issues/74
		parser.Client = &http.Client{
			Transport: &UserAgentTransport{http.DefaultTransport},
		}
		parsed, err := parser.ParseURL(feed.FeedLink)

		if err != nil {
			log.WithFields(log.Fields{
				"feedID":   feed.ID,
				"feedLink": feed.Link,
				"error":    err,
			}).Error("Failed parsing feed")
		} else {
			log.WithFields(log.Fields{
				"feedLink":  feed.FeedLink,
				"feedTitle": feed.Title,
			}).Debug("Successfully parsed feed")

			// todo: set last_scraped on Feed model

			// write each retrieved item to db
			for _, item := range parsed.Items {
				feedItem := models.FeedItem{
					FeedID:      feed.ID,
					Title:       item.Title,
					Description: item.Description,
					Content:     item.Content,
					Link:        item.Link,
					Updated:     utils.TimeToNullTime(item.UpdatedParsed),
					Published:   utils.TimeToNullTime(item.PublishedParsed),
					//Author:      utils.StrToNullStr(item.Author),
					GUID: item.GUID,
				}
				if !feedItem.Exists() {
					log.WithFields(log.Fields{
						"feedLink":     feed.FeedLink,
						"feedTitle":    feed.Title,
						"feedItemGUID": item.GUID,
					}).Debug("Feed item is new, persisting it")
					feedItem.Create()
				} else {
					log.WithFields(log.Fields{
						"feedLink":     feed.FeedLink,
						"feedTitle":    feed.Title,
						"feedItemGUID": item.GUID,
					}).Debug("Feed item already scraped")
				}
			}
		}
	}
}
