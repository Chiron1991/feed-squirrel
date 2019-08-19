package main

import (
	"flag"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var logLevel = flag.String("logLevel", "info", "Set log level")
var httpListen = flag.String("httpListen", "127.0.0.1:3000", "Listen address for HTTP traffic")
var maxConc = flag.Int("maxConc", runtime.NumCPU(), "how many feeds can be scraped concurrently")

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello there!"))
}

func init() {
	// parse command line arguments
	flag.Parse()

	// configure logging
	lvl, err := log.ParseLevel(*logLevel)
	if err != nil {
		lvl = log.InfoLevel
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout) // TODO: make configurable
	log.SetLevel(lvl)
	log.Info("Running with log level " + strings.ToUpper(lvl.String()))

	// determine max concurrency for scraping
	log.Info("Concurrency set to " + strconv.Itoa(*maxConc))
}

func main() {
	// setup HTTP routes
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)

	// spawn HTTP server in goroutine to not block execution
	go func() {
		log.Info("Starting HTTP server on " + *httpListen)
		http.ListenAndServe(*httpListen, r)
	}()

	// create timer that periodically triggers the scraper
	for range time.Tick(time.Second * 10) {
		go ScrapeDueFeeds()
	}
}
