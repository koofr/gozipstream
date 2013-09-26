package gozipstream

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"testing"
	"time"
)

func TestZipStream(t *testing.T) {
	z := NewZipStream()

	buf := new(bytes.Buffer)

	copyDone := make(chan error)

	go func() {
		_, copyErr := io.Copy(buf, z)

		copyDone <- copyErr
	}()

	err := z.AddFile("./zipstream.go", "zipstream.go")

	if err != nil {
		t.Errorf("ZipStream AddFile error: %s", err)
	}

	err = z.Add(nil, "dir/", time.Now())

	if err != nil {
		t.Errorf("ZipStream Add error: %s", err)
	}

	z.End()

	err = <-copyDone

	if err != nil {
		t.Errorf("ZipStream Read error: %s", err)
	}

	b := buf.Bytes()

	reader, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))

	if err != nil {
		t.Errorf("Zip reader error: %s", err)
	}

	if l := len(reader.File); l != 2 {
		t.Errorf("Wrong number of files: %d", l)
	}

	info, _ := os.Stat("./zipstream.go")

	file := reader.File[0]

	if file.Name != "zipstream.go" {
		t.Errorf("Wrong file name: %s", file.Name)
	}

	if file.UncompressedSize64 != uint64(info.Size()) {
		t.Errorf("Wrong file size: %s", file.UncompressedSize64)
	}

	dir := reader.File[1]

	if dir.Name != "dir/" {
		t.Errorf("Wrong dir name: %s", dir.Name)
	}

	if dir.UncompressedSize64 != uint64(0) {
		t.Errorf("Wrong dir size: %s", dir.UncompressedSize64)
	}
}
