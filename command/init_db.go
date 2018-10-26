package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/simplesurance/baur"
	"github.com/simplesurance/baur/log"
	"github.com/simplesurance/baur/storage/postgres"
)

const initDbExample = `
baur init db postgres://postgres@localhost:5432/baur?sslmode=disable
`

const initDbLongHelp = `
Creates the baur tables in a PostgreSQL database.
If no URI is passed, the postgres_uri from the repository config is used.
`

var initDbCmd = &cobra.Command{
	Use:     "db [POSTGRES-URI]",
	Short:   "create baur tables in a PostgreSQL database",
	Example: strings.TrimSpace(initDbExample),
	Long:    initDbLongHelp,
	Run:     initDb,
	Args:    cobra.MaximumNArgs(1),
}

func init() {
	initCmd.AddCommand(initDbCmd)
}

func initDb(cmd *cobra.Command, args []string) {
	var dbURI string

	if len(args) == 0 {
		repo, err := findRepository()
		if err != nil {
			log.Fatalf("could not find '%s' repository config file.\n"+
				"Pass the Postgres URI as argument or run 'baur init repo' first.",
				baur.RepositoryCfgFile)
		}

		dbURI = repo.PSQLURL
	} else {
		dbURI = args[0]
	}

	storageClt, err := postgres.New(dbURI)
	if err != nil {
		log.Fatalln("establishing connection failed:", err.Error())
	}

	err = storageClt.Init()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("database tables created successfully")
}
