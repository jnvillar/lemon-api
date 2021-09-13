package server

import (
	"github.com/spf13/cobra"
)

const (
	defaultLogLevel                = "debug"
	defaultSQLDatabase             = "lemon"
	defaultSQLDatabaseMaxOpenConns = 10
	defaultSQLUser                 = "root"
	defaultSQLPassword             = "your_secret"
)

var (
	logLevelArg             string
	sqlDatabase             string
	sqlDatabaseMaxOpenConns int
	sqlUser                 string
	sqlPassword             string
)

var cmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the Server",
	Long:  `Runs the Server`,
}

func init() {
	cmd.Flags().StringVar(&logLevelArg, "loglevel", defaultLogLevel, "the log level to be considered when logging")
	cmd.Flags().StringVar(&sqlDatabase, "database", defaultSQLDatabase, "SQL database")
	cmd.Flags().StringVar(&sqlUser, "sql_user", defaultSQLUser, "SQL user")
	cmd.Flags().StringVar(&sqlPassword, "sql_password", defaultSQLPassword, "SQL password")
	cmd.Flags().IntVar(&sqlDatabaseMaxOpenConns, "sql_max_open_conns", defaultSQLDatabaseMaxOpenConns, "Max number of open connections to the database. Zero means unlimited.")
}
