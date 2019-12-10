package cfg

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

// Include represents an include configuration file.
type Include struct {
	ID     string `toml:"id" comment:"identifier for the include"`
	Input  Input
	Output Output
}

// ExampleInclude returns an Include struct with exemplary values.
func ExampleInclude() *Include {
	return &Include{
		Input:  exampleInput(),
		Output: exampleOutput(),
	}
}

// IncludeToFile serializes the Include struct to TOML and writes it to filepath.
func (in *Include) IncludeToFile(filepath string) error {
	return toFile(in, filepath, false)
}

// IncludeFromFile deserializes an Include struct from a file.
func IncludeFromFile(path string) (*Include, error) {
	config := Include{}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	config.Output.removeEmptySections()

	return &config, err
}

// Validate validates an Include configuration struct.
func (in *Include) Validate() error {
	if in.ID == "" {
		return errors.New("id can not be empty")
	}

	if err := in.Input.Validate(); err != nil {
		return errors.Wrap(err, "[Input] section contains errors")
	}

	if err := in.Output.Validate(); err != nil {
		return errors.Wrap(err, "[Output] section contains errors")
	}

	return nil
}
