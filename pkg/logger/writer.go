package logger

type CustomWriter struct {
	channel string
	appName string
	tagName string
}

func (w *CustomWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
