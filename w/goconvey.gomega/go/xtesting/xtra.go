package xtesting

// Interface is the interface required of test loggers.
// The os package will invoke the interface's methods to indicate that
// it is inspecting the given environment variables or files.
// Multiple goroutines may call these methods simultaneously.
type TestLogInterface interface {
	Getenv(key string)
	Stat(file string)
	Open(file string)
	Chdir(dir string)
}

func GetTestLog() TestLogInterface {
	return &log
}
