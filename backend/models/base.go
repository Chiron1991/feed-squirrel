package models

import (
	"github.com/Chiron1991/feed-squirrel/backend/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var db *sqlx.DB

type timestampedModel struct {
	ID        int64
	CreatedAt time.Time   `db:"created_at"`
	DeletedAt pq.NullTime `db:"deleted_at" json:"-"`
}

func init() {
	dbConnStr := utils.Getenv("SQUIRREL_DB", "postgres://squirrel:squirrel@db:5432/feedsquirrel?sslmode=disable")

	conn, err := sqlx.Connect("postgres", dbConnStr)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not connect to the database")
	}
	db = conn

	// apply schema migrations
	migrations := &migrate.FileMigrationSource{
		Dir: utils.Getenv("SQUIRREL_MIGRATION_DIR", "models/migrations"),
	}
	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to apply migrations")
	}
	log.Info("Applied " + strconv.Itoa(n) + " schema migrations")
}

// GetDB returns a handle to the database
func GetDB() *sqlx.DB {
	return db
}
