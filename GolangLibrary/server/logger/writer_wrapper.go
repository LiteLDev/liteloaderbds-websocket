package logger

type WriterWrapper struct {
	WriteFunc func(p []byte)
}

func (s WriterWrapper) Write(p []byte) (n int, err error) {
	if p[len(p)-1] == '\n' {
		p = p[:len(p)-1]
	}
	s.WriteFunc(p)
	return len(p), nil
}
