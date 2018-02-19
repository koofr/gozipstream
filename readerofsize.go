package gozipstream

import "io"

type readerOfSize struct {
	size        int64
	alreadyRead int64
}

func newReaderOfSize(size int64) *readerOfSize {
	return &readerOfSize{
		size: size,
	}
}

func (r *readerOfSize) Read(p []byte) (n int, err error) {
	l := len(p)
	l64 := int64(l)
	if r.alreadyRead+l64 >= r.size {
		remaining := r.size - r.alreadyRead
		if remaining < 0 {
			remaining = 0
		}
		r.alreadyRead += remaining
		return int(remaining), io.EOF
	}
	r.alreadyRead += l64
	return l, nil
}
