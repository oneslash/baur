package baur1

import (
	"fmt"

	"github.com/simplesurance/baur/cfg"
	"github.com/simplesurance/baur/fs"
)

// AppCfgFile contains the name of application configuration files
const AppCfgFile = ".app.toml"

type AppLoader struct {
	includedb  *cfg.IncludeDB
	appConfigs map[string]*cfg.App
}

func NewAppLoader(searchDirectories []string, searchDepth int) (*AppLoader, error) {
	a := AppLoader{
		includedb:  cfg.NewIncludeDB(),
		appConfigs: map[string]*cfg.App{},
	}

	for _, searchDir := range searchDirectories {
		if err := fs.DirsExist(searchDir); err != nil {
			return nil, fmt.Errorf("application search directory: %w", err)
		}

		cfgPaths, err := fs.FindFilesInSubDir(searchDir, AppCfgFile, searchDepth)
		if err != nil {
			return nil, fmt.Errorf("discovering application configs failed: %w", err)
		}

		for _, cfgPath := range cfgPaths {
			a.appConfigs[cfgPath] = nil
		}
	}

	return &a, nil
}

func (a *AppLoader) AppConfig(path string) (*App, error) {
}

func (a *AppLoader) Name(path string) (*App, error) {
	// TODO: the paths shoudl be normalized to prevent errors were the same file is referenced via different relative paths
	app, exist := a.appConfigs[path]
}

func (a *AppLoader) All() ([]*App, error) {
}
