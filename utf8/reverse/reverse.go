//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var builder strings.Builder

	builder.Grow(len(input))

	for i := len(input); i > 0; {
		r, size := utf8.DecodeRuneInString(input[i:])

		if r == utf8.RuneError && size == 1 {
			builder.WriteRune(utf8.RuneError)
		} else {
			builder.WriteRune(r)
		}

		i -= size
	}

	return builder.String()
}
