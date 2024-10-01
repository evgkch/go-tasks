//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var builder strings.Builder
	flag := false

	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		if r == utf8.RuneError && size == 1 {
			builder.WriteRune(utf8.RuneError)
			flag = false
		} else if unicode.IsSpace(r) {
			if !flag {
				builder.WriteRune(' ')
				flag = true
			}
		} else {
			builder.WriteRune(r)
			flag = false
		}

		i += size
	}

	return builder.String()
}
