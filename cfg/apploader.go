package cfg

import "fmt"

type AppLoader struct {
	IncludeDB IncludeDB
}

func (a *AppLoader) Load(path string) (*App, error) {
	app, err := AppFromFile(path)
	if err != nil {
		return nil, err
	}

	for _, task := range app.Tasks {
		for _, includeID := range task.Includes {
			include, exist := a.IncludeDB[includeID]
			if !exist {
				return nil, fmt.Errorf("could not find include with id '%s'", includeID)
			}

			task.Merge(include)

			// TODO: store the repository relative cfg path in the include somehow
			//task.Input.Files.Paths = append(task.Input.Files.Paths, include.RelCfgPath)
		}
	}

	return app, nil
}
