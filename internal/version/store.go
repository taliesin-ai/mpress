package version

import (
	"errors"
	"fmt"
	"os"

	"github.com/leaanthony/mpress/internal/config"
	"github.com/leaanthony/mpress/internal/projectfs"
	"github.com/leaanthony/mpress/internal/routes"
)

var errMissingStore = errors.New("version artifacts directory missing")

func validateLabel(label string) error {
	if !labelRE.MatchString(label) || routes.Component(label) != nil {
		return fmt.Errorf("invalid version label %q", label)
	}
	return nil
}

func openVersionStore(project string) (*projectfs.FS, config.Config, error) {
	files, err := projectfs.Open(project)
	if err != nil {
		return nil, config.Config{}, err
	}
	defer files.Close()
	cfg, err := config.LoadWithReadFile(project, files.ReadFile)
	if err == nil {
		err = cfg.SafeVersionPath(project, cfg.ArtifactsPath(project))
	}
	if err != nil {
		return nil, cfg, err
	}
	store, err := files.Sub(cfg.ArtifactsPath(project))
	if os.IsNotExist(err) {
		err = fmt.Errorf("%w: %w", errMissingStore, err)
	}
	return store, cfg, err
}
