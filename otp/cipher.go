package otp

import (
	"io"
)

type reader struct {
	data io.Reader
	prng io.Reader
}

func (s *reader) Read(p []byte) (n int, err error) {
	n, err = s.data.Read(p)
	if n == 0 {
		return n, err
	}

	k := make([]byte, n)
	_, err = s.prng.Read(k)
	if err != nil {
		return 0, err
	}

	for i, x := range k {
		p[i] ^= x
	}

	return n, err // Возвращаем исходную ошибку
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &reader{
		data: r,
		prng: prng,
	}
}

type writer struct {
	data io.Writer
	prng io.Reader
}

func (s *writer) Write(p []byte) (n int, err error) {
	k := make([]byte, len(p))
	_, err = s.prng.Read(k)
	if err != nil {
		return 0, err
	}

	l := make([]byte, len(p))
	for i, x := range k {
		l[i] = x ^ p[i]
	}

	n, err = s.data.Write(l)
	if err != nil {
		return n, err
	}

	return len(p), nil
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &writer{
		data: w,
		prng: prng,
	}
}
