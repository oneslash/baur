package cfg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/simplesurance/baur/fs"
)

// IncludeDB loads and stores include config files
type IncludeDB map[string]*Include

// Load reads and validates all *.toml files in the passed includeDirectories
// as include config files and adds them to the database.
// Directories are searched recursively and symlinks are followed.
func (l IncludeDB) Load(includeDirectory ...string) error {
	walkFunc := func(path string, _ os.FileInfo) error {
		if filepath.Ext(path) != ".toml" {
			return nil
		}

		include, err := IncludeFromFile(path)
		if err != nil {
			return errors.Wrapf(err, "loading include file %q failed", path)
		}

		err = include.Validate()
		if err != nil {
			return errors.Wrapf(err, "validating include config %q failed", path)
		}

		if err := l.add(include); err != nil {
			return err
		}

		return nil
	}

	for _, includeDir := range includeDirectory {
		err := fs.WalkFiles(includeDir, fs.SymlinksAreErrors, walkFunc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l IncludeDB) add(include *Include) error {
	if _, exist := l[include.ID]; exist {
		return fmt.Errorf("multiple includes with id '%s' exist, include ids must be unique", include.ID)
	}

	l[include.ID] = include

	return nil
}
