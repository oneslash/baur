package cfg

type AppLoader struct {
	IncludeDB *IncludeDB
}

func (a *AppLoader) Load(path string) (*App, error) {
	app, err := AppFromFile(path)
	if err != nil {
		return nil, err
	}

	if err := app.Merge(a.IncludeDB); err != nil {
		return nil, err
	}

	return app, nil
}
