package main

import (
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
)

func toDatabase(feed gofeed.Feed) {
	log.WithFields(log.Fields{
		"feedUrl":   feed.Link,
		"feedTitle": feed.Title,
	}).Info("Successfully parsed feed")
}

func handleError(url string, err error) {
	log.WithFields(log.Fields{
		"feedUrl": url,
		"error":   err,
	}).Error("Failed parsing feed")
}

func worker(id int, jobs <-chan string) {
	for url := range jobs {
		log.WithFields(log.Fields{
			"workerId": id,
			"feedUrl":  url,
		}).Debug("Scraping feed")

		parser := gofeed.NewParser()
		feed, err := parser.ParseURL(url)
		if err != nil {
			handleError(url, err)
		} else {
			toDatabase(*feed)
		}
	}
}

func ScrapeDueFeeds() {
	log.Trace("Starting to scrape due feeds")

	feedUrls := []string{
		"https://www.computerbase.de/rss/news.xml",
		"https://blog.fefe.de/rss.xml?html",
		"https://rss.golem.de/rss.php?feed=ATOM1.0",
		"https://hnrss.org/frontpage",
		"https://www.reddit.com/r/pathofexile/search.rss?q=flair:%22GGG%22&sort=new",
		"https://thearmoredpatrol.com/category/world-of-tanks/feed/",
		"https://xkcd.com/rss.xml",
		"http://www.pathofexile.com/news/rss",
		"http://feeds.feedburner.com/blogspot/rkEL",
		"https://www.kernel.org/feeds/kdist.xml",
		"http://feeds.feedburner.com/d0od",
		"http://us2.campaign-archive1.com/feed?u=e2e180baf855ac797ef407fc7&id=9e26887fc5",
	} // TODO: replace with DB lookup

	jobs := make(chan string, 100)

	// spawn workers
	for w := 1; w <= *maxConc; w++ {
		go worker(w, jobs)
	}

	// stuff urls into worker
	for _, url := range feedUrls {
		log.Trace("Stuffing " + url + " into worker")
		jobs <- url
	}
	close(jobs)
	log.Trace("Closed jobs channel")
}
