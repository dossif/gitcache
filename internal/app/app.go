package app

import (
	"context"
	"github.com/dossif/gitcache/internal/config"
	"github.com/dossif/gitcache/internal/services"
	"github.com/dossif/gitcache/pkg/logger"
	"sync"
)

func Run(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, log *logger.Logger, appName string, appVersion string) {
	log.Lg.Info().Msgf("start %v ver %v", appName, appVersion)
	// git cache
	gcSvc := services.NewGitCache(ctx, *cfg, log)
	wg.Add(1)
	go func() { gcSvc.Source(); defer wg.Done() }()
	wg.Add(1)
	go func() { gcSvc.Mirror(); defer wg.Done() }()
	wg.Wait()
	defer log.Lg.Info().Msgf("stop %v", appName)
}
