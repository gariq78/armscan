package usecases

type Logger interface {
	Inf(string, ...interface{})
	Err(error, string, ...interface{})
}
