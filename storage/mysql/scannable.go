package mysql

type Scannable interface {
	Scan(...interface{}) error
}
