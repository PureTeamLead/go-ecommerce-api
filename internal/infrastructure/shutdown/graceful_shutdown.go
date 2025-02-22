package shutdown

import (
	"database/sql"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func Stop(dbConn *sql.DB, logger *zap.Logger, fn func()) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	logger.Info("Gracefully shutdowning the program...")
	logger.Sync()
	dbConn.Close()
	fn()
}
