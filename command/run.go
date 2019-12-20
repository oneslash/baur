package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/simplesurance/baur"
	"github.com/simplesurance/baur/baur1"
	"github.com/simplesurance/baur/fs"
	"github.com/simplesurance/baur/log"
	"github.com/spf13/cobra"
)

// TODO:
// - support specifying only app name, to run all tasks of the app
// - support specifying only task name, to run tasks for all apps with the same name

var runLongHelp = fmt.Sprintf(`
Run Tasks.
If no argument is passed, all tasks in the repository are run,.
By default only tasks with status %s are run.

Tasks-Specifier is in the format:
    <APPLICATION>.<TASK>
    <APPLICATION> or <TASK> can be '*' to match all applications or tasks.

The following Environment Variables are supported:
    %s

  S3 Upload:
    %s
    %s
    %s

  Docker Registry Upload:
    %s
    %s
    %s
    %s
    %s
    %s
`,
	coloredBuildStatus(baur.BuildStatusPending),

	highlight(envVarPSQLURL),

	highlight("AWS_REGION"),
	highlight("AWS_ACCESS_KEY_ID"),
	highlight("AWS_SECRET_ACCESS_KEY"),

	highlight(dockerEnvUsernameVar),
	highlight(dockerEnvPasswordVar),
	highlight("DOCKER_HOST"),
	highlight("DOCKER_API_VERSION"),
	highlight("DOCKER_CERT_PATH"),
	highlight("DOCKER_TLS_VERIFY"))

// TODO: Passing "*" as argument is not nice to use in a shell, without quoting it will expand

const runExample = `
baur run payment-service.build	Run the build task of the payment-service application if it's status is pending.
baur run *.check		Run all check tasks in status pending of all applications
baur run --force --skip-upload	Run all tasks of all application, rerun them if their status is not pending, skip uploading outputs
`

var runCmd = &cobra.Command{
	Use:     "run [<TASK-SPECIFIER>]",
	Short:   "run tasks",
	Long:    strings.TrimSpace(runLongHelp),
	Run:     execRun,
	Example: strings.TrimSpace(buildExampleHelp),
	Args:    cobra.MaximumNArgs(1),
}

var runCmdConf = struct {
	skipRecord bool
	skipUpload bool
	force      bool
}{}

func init() {
	buildCmd.Flags().BoolVarP(&runCmdConf.skipUpload, "skip-upload", "s", false,
		"skip uploading task outputs")
	buildCmd.Flags().BoolVarP(&runCmdConf.skipUpload, "skip-record", "r", false,
		"skip recording the results to the database, --skip-upload must also be passed")
	buildCmd.Flags().BoolVarP(&runCmdConf.force, "force", "f", false,
		"force rebuilding of tasks with status "+baur.BuildStatusExist.String())
	rootCmd.AddCommand(buildCmd)
}

func execRun(cmd *cobra.Command, args []string) {
	if runCmdConf.skipRecord && !runCmdConf.skipUpload {
		log.Fatalln("--skip-upload must be passed when --skip-record is specified")
	}

	if !runCmdConf.skipUpload {
		log.Fatalln("running tasks without --skip-upload is not implemented")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could get current working directory: %v", err)
	}

	repoCfg, err := baur1.FindAndLoadRepositoryConfig(cwd)
	if err != nil {
		log.Fatalln(err)
	}

	repositoryRoot := filepath.Base(repoCfg.FilePath())
	absSearchDirs := fs.AbsPaths(repositoryRoot, repoCfg.Discover.Dirs)

	appLoader := baur1.NewAppLoader(absSearchDirs, repoCfg.Discover.SearchDepth)
	appLoader.All()

	/*
		app := app.App{}
		err := app.RunTask(args[0], runCmdConf.skipRecord, runCmdConf.skipUpload, runCmdConf.force)
		if err != nil {
			log.Fatalln(err)
		}
	*/
}
