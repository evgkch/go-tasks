//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"strings"
)

type lineReader struct {
	reader *bufio.Reader
}

func (lr *lineReader) ReadLine() (string, error) {
	line, err := lr.reader.ReadString('\n')
	line = strings.TrimSuffix(line, "\n")

	if err == io.EOF && len(line) > 0 {
		return line, nil
	}

	if err != nil {
		return "", err
	}

	return line, nil
}

func NewReader(r io.Reader) LineReader {
	return &lineReader{reader: bufio.NewReader(r)}
}

type lineWriter struct {
	writer *bufio.Writer
}

func (lw *lineWriter) Write(l string) error {
	_, err := lw.writer.WriteString(l + "\n")

	if err != nil {
		return err
	}

	return nil
}

func NewWriter(w io.Writer) LineWriter {
	return &lineWriter{writer: bufio.NewWriter(w)}
}

func Merge(w LineWriter, readers ...LineReader) error {
	for i, reader := range readers {

	}

	h := interface{}
	heap.Init(h)

	return nil
}

func Sort(w io.Writer, in ...string) error {
	panic("implement me")
}
