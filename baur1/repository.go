package baur1

import "github.com/simplesurance/baur/cfg"

type Repository struct {
	rootPath  string
	cfg       *cfg.Repository
	appLoader *AppLoader
}

func NewRepository(rootPath string, cfg *cfg.Repository) *Repository {
	return &Repository{
		rootPath:  rootPath,
		cfg:       cfg,
		appLoader: NewAppLoader(rootPath, cfg.Discover.Dirs, cfg.Discover.SearchDepth),
	}
}

func (r *Repository) App(name string) *App {
}

func (r *Repository) App(name string) *App {
}

func (r *Repository) AppByDirectory(path string) *App {
}
