package main

import (
	"github.com/Chiron1991/feed-squirrel/backend/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Chiron1991/feed-squirrel/backend/routes"
	log "github.com/sirupsen/logrus"
)

func init() {
	// configure logging
	lvl, err := log.ParseLevel(utils.Getenv("SQUIRREL_LOGLEVEL", "info"))
	if err != nil {
		lvl = log.InfoLevel
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout) // TODO: make configurable
	log.SetLevel(lvl)
	log.Info("Running with log level " + strings.ToUpper(lvl.String()))
}

func main() {
	// run scraper periodically in goroutine to not block execution
	scraperEnabled, err := strconv.ParseBool(utils.Getenv("SQUIRREL_SCRAPER_ENABLED", "true"))
	if err == nil && scraperEnabled {
		go func() {
			for range time.Tick(time.Second * 10) { // todo: make configurable
				go ScrapeDueFeeds()
			}
		}()
	} else {
		log.Warn("Scraper is disabled")
	}

	// spawn HTTP server
	router := routes.GetRouter() // todo: serve frontend files
	httpListen := utils.Getenv("SQUIRREL_LISTEN", ":3000")
	log.Info("HTTP server listening on " + httpListen)
	log.Fatal(http.ListenAndServe(httpListen, *router))
}
