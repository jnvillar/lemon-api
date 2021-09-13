package server

type Options struct {
	LogLevel                string
	SQLDatabase             string
	SQLDatabaseMaxOpenConns int
	SQLUser                 string
	SQLPassword             string
}
