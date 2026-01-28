//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
	"sort"
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
	if len(readers) == 0 {
		return nil
	}

	h := &MinHeap{}
	heap.Init(h)

	for i, reader := range readers {
		line, err := reader.ReadLine()
		if err == io.EOF {
			continue
		}
		if err != nil {
			return err
		}

		heap.Push(h, &heapItem{line: line, readerIndex: i})
	}

	for h.Len() > 0 {
		item := heap.Pop(h).(*heapItem)

		err := w.Write(item.line)

		if err != nil {
			return err
		}

		line, err := readers[item.readerIndex].ReadLine()
		if err == io.EOF {
			continue
		}
		if err != nil {
			return err
		}

		heap.Push(h, &heapItem{line: line, readerIndex: item.readerIndex})
	}

	if lw, ok := w.(*lineWriter); ok {
		return lw.writer.Flush()
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	if len(in) == 0 {
		return nil
	}

	var files []*os.File
	readers := make([]LineReader, 0, len(in))

	for _, filename := range in {
		// 1. Read
		file, err := os.Open(filename)
		if err != nil {
			return err
		}

		reader := NewReader(file)
		var lines []string
		for {
			line, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				file.Close()
				return err
			}
			lines = append(lines, line)
		}
		file.Close()

		// 2. Sort
		sort.Strings(lines)

		// 3. Write
		file, err = os.Create(filename)
		if err != nil {
			return err
		}

		writer := NewWriter(file)
		for _, line := range lines {
			err := writer.Write(line)
			if err != nil {
				file.Close()
				return err
			}
		}
		if lw, ok := writer.(*lineWriter); ok {
			if err := lw.writer.Flush(); err != nil {
				file.Close()
				return err
			}
		}
		file.Close()

		// 4. Open sorted file for read
		file, err = os.Open(filename)
		if err != nil {
			return err
		}

		files = append(files, file)
		readers = append(readers, NewReader(file))
	}

	defer func() {
		for _, file := range files {
			file.Close()
		}
	}()

	return Merge(NewWriter(w), readers...)
}
