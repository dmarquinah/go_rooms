package routes

import (
	stdlog "log"
	"net/http"
	"os"

	log "github.com/go-kit/log"

	"github.com/dmarquinah/go_rooms/middlewares"
)

func SetupGlobalMiddlewares(router http.Handler) http.Handler {
	return setupLoggingMiddleware(router)
}

func setupLoggingMiddleware(router http.Handler) http.Handler {
	var logger log.Logger
	// Logfmt is a structured, key=val logging format that is easy to read and parse
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	// Direct any attempts to use Go's log package to our structured logger
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	// Log the timestamp (in UTC) and the callsite (file + line number) of the logging
	// call for debugging in the future.
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "loc", log.DefaultCaller)
	// Create an instance of our LoggingMiddleware with our configured logger
	loggingMiddleware := middlewares.LoggingMiddleware(logger)
	return loggingMiddleware(router)
}
