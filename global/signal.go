package global

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForTermSignals() {
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
}
