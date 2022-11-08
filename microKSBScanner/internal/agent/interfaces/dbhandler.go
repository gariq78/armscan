package interfaces

type DbHandler interface {
	Execute(statement string, args ...interface{}) error
	Query(statement string, args ...interface{}) (Rows, error)
}

type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
	Close() error
}
