package utils

import log "github.com/sirupsen/logrus"

// LogSQLError makes a structured log entry about a db query that failed
func LogSQLError(err error, stmt string, args ...interface{}) {
	log.WithFields(log.Fields{
		"error": err,
		"stmt":  stmt,
		"args":  args,
	}).Error("Database statement failed")
}
