package interfaces

type Logger interface {
	Inf(string, ...interface{})
	Err(error, string, ...interface{})
}
