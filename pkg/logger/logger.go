package logger

import (
	"fmt"
	"github.com/grafana/go-gelf/v2/gelf"
	"github.com/rs/zerolog"
	"io"
	"os"
)

type Logger struct {
	Lg *zerolog.Logger
}

func NewLogger(level string, gelfSrv string) (*Logger, error) {
	var logOut io.Writer
	if gelfSrv != "" {
		gelfW, err := gelf.NewTCPWriter(gelfSrv)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to gelf server: %v", err)
		}
		logOut = io.MultiWriter(os.Stdout, gelfW)
	} else {
		logOut = os.Stdout
	}
	lg := Logger{}
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level %v: %w", level, err)
	}
	var newLg zerolog.Logger
	if lvl == zerolog.DebugLevel {
		newLg = zerolog.New(logOut).Level(lvl).With().Caller().Timestamp().Logger()
	} else {
		newLg = zerolog.New(logOut).Level(lvl).With().Logger()
	}
	lg.Lg = &newLg
	return &lg, nil
}

func CheckLevel(lvl string, check zerolog.Level) bool {
	level, _ := zerolog.ParseLevel(lvl)
	if level == check {
		return true
	} else {
		return false
	}
}
