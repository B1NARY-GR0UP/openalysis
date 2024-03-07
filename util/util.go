package util

import (
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// SplitNameWithOwner split nameWithOwner string into owner and name string
// e.g. cloudwego/hertz => cloudwego hertz
func SplitNameWithOwner(s string) (string, string) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func WaitSignal(errC chan error) error {
	signalToNotify := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM}
	if signal.Ignored(syscall.SIGHUP) {
		signalToNotify = signalToNotify[1:]
	}
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, signalToNotify...)
	// block here
	select {
	case sig := <-signalC:
		switch sig {
		case syscall.SIGTERM:
			// force exit
			return errors.New(sig.String())
		case syscall.SIGHUP, syscall.SIGINT:
			// graceful shutdown
			slog.Info("receive signal: ", "signal", sig.String())
			return nil
		}
	case err := <-errC:
		return err
	}
	return nil
}
