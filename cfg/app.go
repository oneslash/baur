package cfg

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

// App stores an application configuration.
type App struct {
	Name  string  `toml:"name" comment:"Name of the application"`
	Tasks []*Task `toml:"Task"`
}

// Task is a task section
type Task struct {
	Name     string   `toml:"name" comment:"Identifies the task, currently the name must be 'build'."`
	Command  string   `toml:"command" commented:"false" comment:"Command that the task executes"`
	Includes []string `toml:"includes" comment:"Repository relative paths to baur include files that the build inherits.\n Valid variables: $ROOT"`
	Input    Input    `toml:"Input" comment:"Specification of task inputs like source files, Makefiles, etc"`
	Output   Output   `toml:"Output" comment:"Specification of task outputs produced by the Task.command"`
}

// Input contains information about task inputs
type Input struct {
	Files         FileInputs    `comment:"Inputs specified by file glob paths"`
	GitFiles      GitFileInputs `comment:"Inputs specified by path, matching only Git tracked files"`
	GolangSources GolangSources `comment:"Inputs specified by directories containing Golang applications"`
}

// GolangSources specifies inputs for Golang Applications
type GolangSources struct {
	Environment []string `toml:"environment" comment:"Environment to use when discovering Golang source files\n This can be environment variables understood by the Golang tools, like GOPATH, GOFLAGS, etc.\n If empty the default Go environment is used.\n Valid variables: $ROOT " commented:"true"`
	Paths       []string `toml:"paths" comment:"Paths to directories containing Golang source files.\n All source files including imported packages are discovered,\n files from Go's stdlib package and testfiles are ignored." commented:"true"`
}

// FileInputs describes a file source
type FileInputs struct {
	Paths []string `toml:"paths" commented:"true" comment:"Relative path to source files,\n supports Golang's Glob syntax (https://golang.org/pkg/path/filepath/#Match) and\n ** to match files recursively\n Valid variables: $ROOT"`
}

// GitFileInputs describes source files that are in the git repository by git
// pathnames
type GitFileInputs struct {
	Paths []string `toml:"paths" commented:"true" comment:"Relative paths to source files.\n Only files tracked by Git that are not in the .gitignore file are matched.\n The same patterns that git ls-files supports can be used.\n Valid variables: $ROOT"`
}

// Output the build output section
type Output struct {
	DockerImage []*DockerImageOutput `comment:"Docker images that are produced by the [Build.command]"`
	File        []*FileOutput        `comment:"Files that are produces by the [Build.command]"`
}

// FileOutput describes where a file artifact should be uploaded to
type FileOutput struct {
	Path     string   `toml:"path" comment:"Path relative to the application directory, valid variables: $APPNAME" commented:"true"`
	FileCopy FileCopy `comment:"Copy the file to a local directory"`
	S3Upload S3Upload `comment:"Upload the file to S3"`
}

// FileCopy describes where a file artifact should be copied to
type FileCopy struct {
	Path string `toml:"path" comment:"Destination directory" commented:"true"`
}

// DockerImageRegistryUpload holds information about where the docker image
// should be uploaded to
type DockerImageRegistryUpload struct {
	Repository string `toml:"repository" comment:"Repository path, format: [<server[:port]>/]<owner>/<repository>:<tag>, valid variables: $APPNAME" commented:"true"`
	Tag        string `toml:"tag" comment:"Tag that is applied to the image, valid variables: $APPNAME, $UUID, $GITCOMMIT" commented:"true"`
}

// S3Upload contains S3 upload information
type S3Upload struct {
	Bucket   string `toml:"bucket" comment:"Bucket name, valid variables: $APPNAME" commented:"true"`
	DestFile string `toml:"dest_file" comment:"Remote File Name, valid variables: $APPNAME, $UUID, $GITCOMMIT" commented:"true"`
}

// DockerImageOutput describes where a docker container is uploaded to
type DockerImageOutput struct {
	IDFile         string                    `toml:"idfile" comment:"Path to a file that is created by [Build.Command] and contains the image ID of the produced image (docker build --iidfile), valid variables: $APPNAME" commented:"true"`
	RegistryUpload DockerImageRegistryUpload `comment:"Registry repository the image is uploaded to"`
}

func exampleBuildInput() Input {
	return Input{
		Files: FileInputs{
			Paths: []string{"dbmigrations/*.sql"},
		},
		GitFiles: GitFileInputs{
			Paths: []string{"Makefile"},
		},
		GolangSources: GolangSources{
			Paths:       []string{"."},
			Environment: []string{"GOFLAGS=-mod=vendor", "GO111MODULE=on"},
		},
	}
}

func exampleBuildOutput() Output {
	return Output{
		File: []*FileOutput{
			{
				Path: "dist/$APPNAME.tar.xz",
				S3Upload: S3Upload{
					Bucket:   "go-artifacts/",
					DestFile: "$APPNAME-$GITCOMMIT.tar.xz",
				},
				FileCopy: FileCopy{
					Path: "/mnt/fileserver/build_artifacts/$APPNAME-$GITCOMMIT.tar.xz",
				},
			},
		},
		DockerImage: []*DockerImageOutput{
			{
				IDFile: fmt.Sprintf("$APPNAME-container.id"),
				RegistryUpload: DockerImageRegistryUpload{
					Repository: "my-company/$APPNAME",
					Tag:        "$GITCOMMIT",
				},
			},
		},
	}
}

// ExampleApp returns an exemplary app cfg struct with the name set to the given value
func ExampleApp(name string) *App {
	return &App{
		Name: name,

		Tasks: []*Task{
			&Task{
				Name:    "build",
				Command: "make dist",
				Input:   exampleBuildInput(),
				Output:  exampleBuildOutput(),
			},
		},
	}
}

// AppFromFile reads a application configuration file and returns it.
// If the buildCmd is not set in the App configuration it's set to
// defaultBuild.Command
func AppFromFile(path string) (*App, error) {
	config := App{}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	for _, task := range config.Tasks {
		removeEmptySections(&task.Output)
	}

	return &config, err
}

// removeEmptySections removes elements from slices of the that are empty.
// This is a workaround for https://github.com/pelletier/go-toml/issues/216
// It prevents that slices are commented in created Example configurations.
// To prevent that we have empty elements in the slice that we process later and
// validate, remove them from the config
func removeEmptySections(buildOutput *Output) {
	fileOutputs := make([]*FileOutput, 0, len(buildOutput.File))
	dockerImageOutputs := make([]*DockerImageOutput, 0, len(buildOutput.DockerImage))

	for _, f := range buildOutput.File {
		if f.IsEmpty() {
			continue
		}

		fileOutputs = append(fileOutputs, f)
	}

	for _, d := range buildOutput.DockerImage {
		if d.IsEmpty() {
			continue
		}

		dockerImageOutputs = append(dockerImageOutputs, d)
	}

	buildOutput.File = fileOutputs
	buildOutput.DockerImage = dockerImageOutputs
}

// ToFile writes an exemplary Application configuration file to
// filepath. The name setting is set to appName
func (a *App) ToFile(filepath string) error {
	return toFile(a, filepath, false)
}

// Validate validates a App configuration
func (a *App) Validate() error {
	if len(a.Name) == 0 {
		return errors.New("name parameter can not be empty")
	}

	if len(a.Tasks) != 1 {
		return errors.New("The Tasks section must define exactly 1 Task")
	}

	duplMap := make(map[string]struct{}, len(a.Tasks))

	for _, task := range a.Tasks {
		_, exist := duplMap[task.Name]
		if exist {
			return fmt.Errorf("Tasks section contains multiple tasks with the name '%s', task names must be unique", task.Name)
		}
		duplMap[task.Name] = struct{}{}

		err := task.Validate()
		if err != nil {
			if task.Name == "" {
				return errors.Wrap(err, "Task section contains errors")
			}

			return errors.Wrap(err, fmt.Sprintf("Task section %s contains errors", task.Name))
		}
	}

	return nil
}

// Validate validates the build section
func (t *Task) Validate() error {
	if len(t.Command) == 0 {
		return nil
	}

	// TODO: change it to check for an invalid name when we support multiple tasks
	if t.Name != "build" {
		return fmt.Errorf("invalid task name '%s', task name must be 'build'", t.Name)
	}

	if err := t.Input.Validate(); err != nil {
		return errors.Wrap(err, "Input section contains errors")
	}

	if err := t.Output.Validate(); err != nil {
		return errors.Wrap(err, "Output section contains errors")
	}

	return nil
}

// Validate validates the BuildInput section
func (i *Input) Validate() error {
	if err := i.Files.Validate(); err != nil {
		return errors.Wrap(err, "Files")
	}

	if err := i.GolangSources.Validate(); err != nil {
		return errors.Wrap(err, "GolangSources")
	}

	// TODO: add validation for gitfiles section

	return nil
}

// Validate validates the GolangSources section
func (g *GolangSources) Validate() error {
	if len(g.Environment) != 0 && len(g.Paths) == 0 {
		return errors.New("path must be set if environment is set")
	}

	for _, p := range g.Paths {
		if len(p) == 0 {
			return errors.New("a path can not be empty")
		}
	}

	return nil
}

// Validate validates the BuildOutput section
func (o *Output) Validate() error {
	for _, f := range o.File {
		if err := f.Validate(); err != nil {
			return errors.Wrap(err, "File")
		}
	}

	for _, d := range o.DockerImage {
		if err := d.Validate(); err != nil {
			return errors.Wrap(err, "DockerImage")
		}
	}

	return nil
}

// IsEmpty returns true if FileCopy is empty
func (f *FileCopy) IsEmpty() bool {
	return len(f.Path) == 0
}

// IsEmpty returns true if FileOutput is empty
func (f *FileOutput) IsEmpty() bool {
	return f.FileCopy.IsEmpty() && f.S3Upload.IsEmpty()
}

// IsEmpty returns true if S3Upload is empty
func (s *S3Upload) IsEmpty() bool {
	return len(s.Bucket) == 0 && len(s.DestFile) == 0
}

// Validate validates a [[Build.Output.File]] section
func (f *FileOutput) Validate() error {
	if len(f.Path) == 0 {
		return errors.New("path can not be unset or empty")
	}

	return f.S3Upload.Validate()
}

//IsEmpty returns true if the struct is empty
func (d *DockerImageRegistryUpload) IsEmpty() bool {
	return len(d.Repository) == 0 && len(d.Tag) == 0
}

// IsEmpty returns true if DockerImageOutput is empty
func (d *DockerImageOutput) IsEmpty() bool {
	return len(d.IDFile) == 0 && d.RegistryUpload.IsEmpty()

}

// Validate validates a [[Build.Output.File]] section
func (s *S3Upload) Validate() error {
	if s.IsEmpty() {
		return nil
	}

	if len(s.DestFile) == 0 {
		return errors.New("destfile parameter can not be unset or empty")
	}

	if len(s.Bucket) == 0 {
		return errors.New("bucket parameter can not be unset or empty")
	}

	return nil
}

// Validate validates its content
func (d *DockerImageOutput) Validate() error {
	if len(d.IDFile) == 0 {
		return errors.New("idfile parameter can not be unset or empty")
	}

	if err := d.RegistryUpload.Validate(); err != nil {
		return errors.Wrap(err, "") // TODO add section name to error msg
	}

	return nil
}

// Validate validates its content
func (d *DockerImageRegistryUpload) Validate() error {
	if len(d.Repository) == 0 {
		return errors.New("repository parameter can not be unset or empty")
	}

	if len(d.Tag) == 0 {
		return errors.New("tag parameter can not be unset or empty")
	}

	return nil
}

// Validate validates a [[Sources.Files]] section
func (f *FileInputs) Validate() error {
	for _, path := range f.Paths {
		if len(path) == 0 {
			return errors.New("path can not be empty")
		}
		if strings.Count(path, "**") > 1 {
			return errors.New("'**' can only appear one time in a path")
		}
	}

	return nil
}
