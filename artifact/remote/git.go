package remote

import (
	"io/ioutil"
	"net/url"
	"os"

	"golang.org/x/xerrors"

	git "github.com/go-git/go-git/v5"

	"github.com/sf9133/fanal/analyzer"
	"github.com/sf9133/fanal/analyzer/config"
	"github.com/sf9133/fanal/artifact"
	"github.com/sf9133/fanal/artifact/local"
	"github.com/sf9133/fanal/cache"
)

func NewArtifact(rawurl string, c cache.ArtifactCache, disabled []analyzer.Type, opt config.ScannerOption) (
	artifact.Artifact, func(), error) {
	cleanup := func() {}

	u, err := newURL(rawurl)
	if err != nil {
		return nil, cleanup, err
	}

	tmpDir, err := ioutil.TempDir("", "fanal-remote")
	if err != nil {
		return nil, cleanup, err
	}

	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:      u.String(),
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		return nil, cleanup, xerrors.Errorf("git error: %w", err)
	}

	cleanup = func() {
		_ = os.RemoveAll(tmpDir)
	}

	// JAR/WAR/EAR doesn't need to be analyzed in git repositories.
	disabled = append(disabled, analyzer.TypeJar)

	art, err := local.NewArtifact(tmpDir, c, disabled, opt)
	if err != nil {
		return nil, cleanup, xerrors.Errorf("fs artifact: %w", err)
	}
	return art, cleanup, nil
}

func newURL(rawurl string) (*url.URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	// "https://" can be omitted
	// e.g. github.com/sf9133/fanal
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	return u, nil
}
