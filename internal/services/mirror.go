package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/sosedoff/gitkit"
	"log"
	"net/http"
)

func (gc *GC) Mirror() {
	gc.Log.Lg.Info().Str("repo", gc.Repo).Msg("start mirror repo")
	defer gc.Log.Lg.Info().Str("repo", gc.Repo).Msg("stop mirror repo")
	// configure git mirror
	mirror := gitkit.New(gitkit.Config{
		Dir:        gc.Cfg.GitSource.Path,
		AutoCreate: false,
		Auth:       gc.Cfg.GitMirror.Auth,
	})
	if gc.Cfg.GitMirror.Auth == true {
		mirror.AuthFunc = func(cred gitkit.Credential, req *gitkit.Request) (bool, error) {
			log.Println("user auth request for repo:", cred.Username, cred.Password, req.RepoName)
			if cred.Username == gc.Cfg.GitMirror.BasicAuth.Username.Get() &&
				cred.Password == gc.Cfg.GitMirror.BasicAuth.Password.Get() {
				return true, nil
			} else {
				return false, fmt.Errorf("repo basic-auth failed")
			}
		}
	}
	err := mirror.Setup()
	if err != nil {
		gc.Log.Lg.Error().Msgf("failed to setup mirror repo: %v", err)
	}
	// start http server
	httpServer := &http.Server{
		Handler:  mirror,
		Addr:     gc.Cfg.Listen,
		ErrorLog: nil,
	}
	go func() {
		<-gc.Ctx.Done()
		_ = httpServer.Shutdown(context.Background())
	}()
	err = httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) && err != nil {
		gc.Log.Lg.Fatal().Msgf("failed to start httpserver: %v", err)
	}
}
