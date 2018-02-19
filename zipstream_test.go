package gozipstream

import (
	"archive/zip"
	"bytes"
	"fmt"
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
		t.Errorf("Wrong file size: %d", file.UncompressedSize64)
	}

	dir := reader.File[1]

	if dir.Name != "dir/" {
		t.Errorf("Wrong dir name: %s", dir.Name)
	}

	if dir.UncompressedSize64 != uint64(0) {
		t.Errorf("Wrong dir size: %d", dir.UncompressedSize64)
	}
}

func TestZipStreamAddSize(t *testing.T) {
	data := make([]byte, 1*1024*1024+42)

	z := NewZipStream()

	go func() {
		err := z.AddSize(int64(len(data)), "test.bin", time.Now())
		if err != nil {
			t.Errorf("ZipStream AddFile error: %s", err)
		}

		err = z.End()
		if err != nil {
			t.Errorf("ZipStream Read error: %s", err)
		}
	}()

	estimatedTotalSize, err := z.TotalSize()
	if err != nil {
		t.Errorf("ZipStream TotalSize error: %s", err)
	}
	if estimatedTotalSize != 1048748 {
		t.Errorf("expected estimatedTotalSize to be 1048748 but is %d", estimatedTotalSize)
	}

	z = NewZipStream()

	go func() {
		err := z.Add(bytes.NewBuffer(data), "test.bin", time.Now())
		if err != nil {
			t.Errorf("ZipStream AddFile error: %s", err)
		}

		err = z.End()
		if err != nil {
			t.Errorf("ZipStream Read error: %s", err)
		}
	}()

	totalSize, err := z.TotalSize()
	if err != nil {
		t.Errorf("ZipStream TotalSize error: %s", err)
	}
	if totalSize != 1048748 {
		t.Errorf("expected totalSize to be 1048748 but is %d", totalSize)
	}
}

func TestZipStreamAddSizeBenchmark(t *testing.T) {
	var size int64 = 10 * 1024 * 1024 * 1024 // 10 GB

	z := NewZipStream()

	go func() {
		err := z.AddSize(size, "test.bin", time.Now())
		if err != nil {
			t.Errorf("ZipStream AddFile error: %s", err)
		}

		err = z.End()
		if err != nil {
			t.Errorf("ZipStream Read error: %s", err)
		}
	}()

	start := time.Now()

	estimatedTotalSize, err := z.TotalSize()
	if err != nil {
		t.Errorf("ZipStream TotalSize error: %s", err)
	}
	if estimatedTotalSize != 10737418482 {
		t.Errorf("expected estimatedTotalSize to be 10737418482 but is %d", estimatedTotalSize)
	}

	fmt.Printf("%s to copy 10 GB\n", time.Now().Sub(start))
}
