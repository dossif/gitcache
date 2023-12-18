package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dossif/gitcache/internal/app"
	"github.com/dossif/gitcache/internal/config"
	"github.com/dossif/gitcache/pkg/logger"
	"github.com/rs/zerolog"
	"log"
	"os"
	"os/signal"
	"sync"
)

const (
	appName      = "gitcache"
	configPrefix = appName
)

var appVersion = "0.0.0"

func main() {
	fVer := flag.Bool("v", false, "print app version")
	fHlp := flag.Bool("h", false, "print app help")
	fEnv := flag.Bool("e", false, "print app env")
	flag.Parse()
	if *fHlp == true {
		fmt.Println("app params:")
		flag.PrintDefaults()
		os.Exit(128)
	} else if *fVer {
		fmt.Println(fmt.Sprintf("%v v%v", appName, appVersion))
		os.Exit(0)
	} else if *fEnv == true {
		config.PrintUsage(configPrefix)
	}
	wg := new(sync.WaitGroup)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	cfg, err := config.NewConfig(configPrefix)
	if logger.CheckLevel(cfg.LogLevel, zerolog.DebugLevel) == true {
		config.DebugConfig(cfg)
	}
	if err != nil {
		log.Printf("failed to create config: %v", err)
		config.PrintUsage(configPrefix)
	}
	lg, err := logger.NewLogger(cfg.LogLevel, cfg.LogGelf)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	app.Run(ctx, wg, cfg, lg, appName, appVersion)
}
