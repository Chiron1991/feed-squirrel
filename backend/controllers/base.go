package controllers

import (
	"encoding/json"
	"github.com/Chiron1991/feed-squirrel/backend/models"
	"github.com/Chiron1991/feed-squirrel/backend/utils"
	"github.com/go-chi/chi/middleware"
	migrate "github.com/rubenv/sql-migrate"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func respondJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.WithFields(log.Fields{
			"data":      data,
			"requestID": middleware.GetReqID(r.Context()),
		}).Error("Failed to encode JSON for response")
		http.Error(w, http.StatusText(422), 422)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(encoded)
}

// todo: remove later
func LoadDevFixtures(w http.ResponseWriter, r *http.Request) {
	fefe := models.Feed{
		Title:    "Fefe's Blog",
		FeedLink: "https://blog.fefe.de/rss.xml?html",
	}
	fefe.Create()

	cb := models.Feed{
		Title:    "ComputerBase",
		FeedLink: "https://www.computerbase.de/rss/news.xml",
	}
	cb.Create()

	golem := models.Feed{
		Title:    "Golem",
		FeedLink: "https://rss.golem.de/rss.php?feed=ATOM1.0",
	}
	golem.Create()

	hn := models.Feed{
		Title:    "Hacker News: Front Page",
		FeedLink: "https://hnrss.org/frontpage",
	}
	hn.Create()

	kernel := models.Feed{
		Title:    "Linux Kernel Releases",
		FeedLink: "https://www.kernel.org/feeds/kdist.xml",
	}
	kernel.Create()

	omg := models.Feed{
		Title:    "OMG! Ubuntu",
		FeedLink: "http://feeds.feedburner.com/d0od",
	}
	omg.Create()

	respondJSON(w, r, map[string]string{"loaded": "ok"})
}

// todo: remove later
func ClearDB(w http.ResponseWriter, r *http.Request) {
	db := models.GetDB()
	db.MustExec("DROP TABLE feed_items; DROP TABLE feeds; DROP TABLE gorp_migrations;")

	migrations := &migrate.FileMigrationSource{
		Dir: utils.Getenv("SQUIRREL_MIGRATION_DIR", "models/migrations"),
	}
	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to apply migrations")
	}

	respondJSON(w, r, map[string]string{"reset": "ok"})
}
