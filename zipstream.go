package gozipstream

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/koofr/gozipstream/zip"
)

type ZipStream struct {
	reader *io.PipeReader
	writer *io.PipeWriter
	zip    *zip.Writer

	filesSize int64
}

func NewZipStream() (z *ZipStream) {
	reader, writer := io.Pipe()

	z = &ZipStream{
		reader: reader,
		writer: writer,
		zip:    zip.NewWriter(writer),

		filesSize: 0,
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

func (z *ZipStream) AddSize(size int64, name string, mtime time.Time) (err error) {
	return z.add(nil, size, name, mtime)
}

func (z *ZipStream) Add(reader io.Reader, name string, mtime time.Time) (err error) {
	return z.add(reader, 0, name, mtime)
}

func (z *ZipStream) add(reader io.Reader, size int64, name string, mtime time.Time) (err error) {
	isDir := strings.HasSuffix(name, "/")

	header := &zip.FileHeader{
		Name:           name,
		CreatorVersion: 0x0300, // unix
		Flags:          0x800,  // utf8
		Method:         zip.Store,
		Modified:       mtime,
	}

	if isDir {
		header.SetMode(0755)
	} else {
		header.SetMode(0644)
	}

	f, err := z.zip.CreateHeader(header)

	if err != nil {
		return err
	}

	if !isDir {
		if reader != nil {
			_, err = io.Copy(f, reader)
			if err != nil {
				return err
			}
		} else {
			z.filesSize += size

			err = f.(interface {
				AddSize(int64) error
			}).AddSize(size)
			if err != nil {
				return err
			}
		}
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

func (z *ZipStream) TotalSize() (totalSize int64, err error) {
	buf := make([]byte, 64*1024)

	for {
		n, err := z.Read(buf)
		totalSize += int64(n)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}

	err = z.Close()
	if err != nil {
		return 0, err
	}

	totalSize += z.filesSize

	return totalSize, nil
}
