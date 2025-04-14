package main

import (
	"errors"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fInfo, err := os.Stat(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	if offset >= fInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	fSource, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fSource.Close()

	tTarget, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer tTarget.Close()

	if limit == 0 {
		limit = fInfo.Size()
	}
	if offset > 0 {
		_, err = fSource.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	reader := io.LimitReader(fSource, limit)

	if ShowProgress {
		err = CopyWithProgress(tTarget, reader, fInfo.Size())
	} else {
		_, err = io.Copy(tTarget, reader)
	}

	if err != nil {
		return err
	}

	return nil
}

func CopyWithProgress(writer io.Writer, reader io.Reader, total int64) error {
	bar := progressbar.NewOptions64(total)
	_, err := io.Copy(io.MultiWriter(writer, bar), reader)

	return err
}
