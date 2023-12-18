package services

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
)

func (gc *GC) updateGitSourceIter() {
	pe := pathExists(gc.Repo)
	if pe == false {
		gc.Log.Lg.Info().Msgf("clone git repo")
		err := gc.cloneRepo()
		if err != nil {
			gc.Log.Lg.Error().Msgf("failed to clone git repo: %v", err)
			return
		}
	} else {
		gc.Log.Lg.Info().Msgf("git repo already cloned")
	}
	gc.Log.Lg.Info().Msgf("fetch git changes")
	err := gc.fetchRepo()
	if err != nil {
		gc.Log.Lg.Error().Msgf("failed to fetch git repo: %v", err)
		return
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (gc *GC) cloneRepo() error {
	repoSrc := gc.Cfg.GitSource.Url
	var repoAuth transport.AuthMethod
	if gc.Cfg.GitSource.Auth == true {
		repoAuth = &http.BasicAuth{
			Username: gc.Cfg.GitSource.BasicAuth.Username.Get(),
			Password: gc.Cfg.GitSource.BasicAuth.Password.Get(),
		}
	}
	gco := git.CloneOptions{
		URL:               repoSrc,
		Auth:              repoAuth,
		RemoteName:        "",
		ReferenceName:     "",
		SingleBranch:      false,
		Mirror:            true,
		NoCheckout:        false,
		Depth:             0,
		RecurseSubmodules: 0,
		Progress:          gc.Log.Lg,
		Tags:              0,
		InsecureSkipTLS:   false,
		CABundle:          nil,
		ProxyOptions:      transport.ProxyOptions{},
	}
	// clone repo
	_, err := git.PlainClone(gc.Repo, true, &gco)
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		gc.Log.Lg.Info().Str("repo", gc.Repo).Str("src", repoSrc).Msg("git repo already cloned")
	} else if err != nil {
		return fmt.Errorf("failed to clone git repo %v: %v", repoSrc, err)
	}
	return nil
}

func (gc *GC) fetchRepo() error {
	repoSrc := gc.Cfg.GitSource.Url
	repo, err := git.PlainOpen(gc.Repo)
	if err != nil {
		return fmt.Errorf("failed to open git repo %v: %v", gc.Repo, err)
	}
	var repoAuth transport.AuthMethod
	if gc.Cfg.GitSource.Auth == true {
		repoAuth = &http.BasicAuth{
			Username: gc.Cfg.GitSource.BasicAuth.Username.Get(),
			Password: gc.Cfg.GitSource.BasicAuth.Password.Get(),
		}
	}
	gfo := git.FetchOptions{
		RemoteName:      "",
		RemoteURL:       "",
		RefSpecs:        nil,
		Depth:           0,
		Auth:            repoAuth,
		Progress:        gc.Log.Lg,
		Tags:            0,
		Force:           false,
		InsecureSkipTLS: false,
		CABundle:        nil,
		ProxyOptions:    transport.ProxyOptions{},
	}
	err = repo.Fetch(&gfo)
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		gc.Log.Lg.Info().Str("repo", gc.Repo).Str("src", repoSrc).Msg("git repo already up-to-date")
	} else if err != nil {
		return fmt.Errorf("failed to fetch changes from git repo %v: %v", repoSrc, err)
	}
	return nil
}
