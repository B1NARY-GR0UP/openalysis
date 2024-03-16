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

func AssembleDSN(host, port, user, password, database string) string {
	var sb strings.Builder
	sb.WriteString(user)
	sb.WriteString(":")
	sb.WriteString(password)
	sb.WriteString("@tcp(")
	sb.WriteString(host)
	sb.WriteString(":")
	sb.WriteString(port)
	sb.WriteString(")/")
	sb.WriteString(database)
	return sb.String()
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

func IsEmptySlice[T any](slice []T) bool {
	if slice == nil || len(slice) == 0 {
		return true
	}
	return false
}
