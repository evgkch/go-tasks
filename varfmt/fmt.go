package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	var result strings.Builder
	result.Grow(len(format)) // предварительная аллокация

	autoIndex := 0 // счетчик для {}

	for i := 0; i < len(format); i++ {
		if format[i] == '{' {
			// Ищем закрывающую скобку
			closePos := -1
			for j := i + 1; j < len(format); j++ {
				if format[j] == '}' {
					closePos = j
					break
				}
			}

			if closePos == -1 {
				// Нет закрывающей скобки - просто добавляем '{'
				result.WriteByte('{')
				continue
			}

			content := format[i+1 : closePos]

			if content == "" {
				// {} - используем автоиндекс
				if autoIndex < len(args) {
					result.WriteString(fmt.Sprint(args[autoIndex]))
				}
				autoIndex++
			} else {
				// {number} - используем явный индекс
				index, err := strconv.Atoi(content)
				if err != nil {
					// Не число - добавляем как есть
					result.WriteByte('{')
					result.WriteString(content)
					result.WriteByte('}')
				} else {
					if index >= 0 && index < len(args) {
						result.WriteString(fmt.Sprint(args[index]))
					}
					autoIndex++ // ← ВАЖНО! {number} тоже увеличивает счетчик позиций
				}
			}

			i = closePos // ← Переходим к символу после '}'
		} else {
			result.WriteByte(format[i])
		}
	}

	return result.String()
}
