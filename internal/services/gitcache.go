package services

import (
	"context"
	"github.com/dossif/gitcache/internal/config"
	"github.com/dossif/gitcache/pkg/logger"
	"path"
	"time"
)

type GC struct {
	Ctx  context.Context
	Log  logger.Logger
	Cfg  config.Config
	Repo string
}

func NewGitCache(ctx context.Context, cfg config.Config, log *logger.Logger) *GC {
	repoSrc := cfg.GitSource.Url
	repoName := path.Base(repoSrc)
	repoPath := path.Join(cfg.GitSource.Path, repoName)
	return &GC{
		Ctx:  ctx,
		Log:  *log,
		Cfg:  cfg,
		Repo: repoPath,
	}
}

func (gc *GC) Source() {
	gc.Log.Lg.Info().Msg("start git source")
	defer gc.Log.Lg.Info().Msg("stop git source")
	go func() {
		for {
			gc.Log.Lg.Info().Msgf("check and wait %v", gc.Cfg.GitSource.CheckInterval.String())
			gc.updateGitSourceIter()
			time.Sleep(gc.Cfg.GitSource.CheckInterval)
		}
	}()
	<-gc.Ctx.Done()
}
