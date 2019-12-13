package cfg

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

type ValidationError struct {
	ElementPath []string
	Message     string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", strings.Join(v.ElementPath, "."), v.Message)
}

// PrependValidationErrorpath if the passed error has the *ValidationError, the
// passed path is prepended to it's ElementPath field.
// The function returns the passed validationError.
func PrependValidationErrorPath(validationError error, path ...string) error {
	valError, ok := validationError.(*ValidationError)
	if ok {
		valError.ElementPath = append(path, valError.ElementPath...)
	}

	return validationError
}

// toFile serializes a struct to TOML format and writes it to a file.
func toFile(data interface{}, filepath string, overwrite bool) error {
	var openFlags int

	if overwrite {
		openFlags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	} else {
		openFlags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	}

	f, err := os.OpenFile(filepath, openFlags, 0640)
	if err != nil {
		return err
	}

	encoder := toml.NewEncoder(f)
	encoder.Order(toml.OrderPreserve)
	err = encoder.Encode(data)
	if err != nil {
		f.Close()
		return err
	}

	err = f.Close()
	if err != nil {
		return errors.Wrap(err, "closing file failed")
	}

	return err
}
