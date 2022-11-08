package usecases

type Logger interface {
	Err(error, string, ...interface{})
	Inf(string, ...interface{})
}
