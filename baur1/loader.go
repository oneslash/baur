package baur1

import (
	"fmt"

	"github.com/simplesurance/baur/cfg"
	"github.com/simplesurance/baur/fs"
)

const RepositoryCfgFile = ".baur.toml"

// FindRepositoryConfig searches for the RepositoryCfgFile. The search starts in
// the passed directory and traverses the parent directories down to '/'. The
// absolute path to the first found RepositoryCfgFile is returned.
// If the config file is not found os.ErrNotExist is returned.
func FindRepositoryConfig(startPath string) (string, error) {
	cfgPath, err := fs.FindFileInParentDirs(startPath, RepositoryCfgFile)
	if err != nil {
		return "", err
	}

	return cfgPath, nil
}

func FindAndLoadRepositoryConfig(startPath string) (*cfg.Repository, error) {
	cfgPath, err := FindRepositoryConfig(startPath)
	if err != nil {
		return nil, err
	}

	cfg, err := cfg.RepositoryFromFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("loading repository config %q failed: %w", cfgPath, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating repository config %q failed: %w", cfgPath, err)
	}

	return cfg, nil
}
