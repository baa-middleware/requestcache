package requestcache

import "io"

type ghostWriter struct {
	Writer  io.Writer
	Content []byte
}

func (w *ghostWriter) Write(p []byte) (n int, err error) {
	w.Content = append(w.Content, p...)
	n, err = w.Writer.Write(p)
	return
}
