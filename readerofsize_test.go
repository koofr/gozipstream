package gozipstream

import (
	"io"
	"testing"
)

func TestReaderOfSize(t *testing.T) {
	r := newReaderOfSize(10)

	n, err := r.Read(make([]byte, 5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected n to be %d but is %d", 5, n)
	}
	n, err = r.Read(make([]byte, 4))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 4 {
		t.Errorf("expected n to be %d but is %d", 4, n)
	}
	n, err = r.Read(make([]byte, 3))
	if err != io.EOF {
		t.Errorf("expected err to be io.EOF but is %v", err)
	}
	if n != 1 {
		t.Errorf("expected n to be %d but is %d", 1, n)
	}
	n, err = r.Read(make([]byte, 3))
	if err != io.EOF {
		t.Errorf("expected err to be io.EOF but is %v", err)
	}
	if n != 0 {
		t.Errorf("expected n to be %d but is %d", 0, n)
	}

	r = newReaderOfSize(10)

	n, err = r.Read(make([]byte, 10))
	if err != io.EOF {
		t.Errorf("expected err to be io.EOF but is %v", err)
	}
	if n != 10 {
		t.Errorf("expected n to be %d but is %d", 10, n)
	}
	n, err = r.Read(make([]byte, 3))
	if err != io.EOF {
		t.Errorf("expected err to be io.EOF but is %v", err)
	}
	if n != 0 {
		t.Errorf("expected n to be %d but is %d", 0, n)
	}

	r = newReaderOfSize(10)

	n, err = r.Read(make([]byte, 11))
	if err != io.EOF {
		t.Errorf("expected err to be io.EOF but is %v", err)
	}
	if n != 10 {
		t.Errorf("expected n to be %d but is %d", 10, n)
	}
	n, err = r.Read(make([]byte, 3))
	if err != io.EOF {
		t.Errorf("expected err to be io.EOF but is %v", err)
	}
	if n != 0 {
		t.Errorf("expected n to be %d but is %d", 0, n)
	}
}
