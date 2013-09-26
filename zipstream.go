package gozipstream

import (
	"archive/zip"
	"io"
	"os"
	"strings"
	"time"
)

type ZipStream struct {
	reader *io.PipeReader
	writer *io.PipeWriter
	zip    *zip.Writer
}

func NewZipStream() (z *ZipStream) {
	reader, writer := io.Pipe()

	z = &ZipStream{
		reader: reader,
		writer: writer,
		zip:    zip.NewWriter(writer),
	}

	return
}

func (z *ZipStream) AddFile(fullPath string, relPath string) (err error) {
	info, err := os.Stat(fullPath)

	if err != nil {
		return
	}

	file, err := os.Open(fullPath)

	if err != nil {
		return
	}

	defer file.Close()

	return z.Add(file, relPath, info.ModTime())
}

func (z *ZipStream) Add(reader io.Reader, name string, mtime time.Time) (err error) {
	isDir := strings.HasSuffix(name, "/")

	header := &zip.FileHeader{
		Name:           name,
		Method:         zip.Store,
		CreatorVersion: 0x0300, // unix
		Flags:          0x800,  // utf8
	}

	header.SetModTime(mtime)

	if isDir {
		header.SetMode(0755)
	} else {
		header.SetMode(0644)
	}

	f, err := z.zip.CreateHeader(header)

	if err != nil {
		return
	}

	if !isDir {
		_, err = io.Copy(f, reader)
	}

	return
}

func (z *ZipStream) End() (err error) {
	defer z.writer.Close()

	return z.zip.Close()
}

func (z *ZipStream) Error(err error) error {
	return z.writer.CloseWithError(err)
}

func (z *ZipStream) Read(p []byte) (n int, err error) {
	return z.reader.Read(p)
}

func (z *ZipStream) Close() error {
	return z.reader.Close()
}
