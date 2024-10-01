//go:build !solution

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Evaluator struct {
	stack Stack
	words map[string][]string
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		words: make(map[string][]string),
	}
}

func (e *Evaluator) Process(row string) ([]int, error) {
	parts := strings.Fields(strings.ToLower(row))

	err := e.processTokens(parts)
	if err != nil {
		return nil, err
	}

	return e.stack.items, nil
}

func (e *Evaluator) isBuiltin(word string) bool {
	switch word {
	case "+", "-", "*", "/", "dup", "drop", "swap", "over":
		return true
	}
	return false
}

func (e *Evaluator) processTokens(tokens []string) error {
	i := 0
	for i < len(tokens) {
		word := tokens[i]

		if word == ":" {
			if i+2 >= len(tokens) {
				return fmt.Errorf("invalid definition")
			}

			endIdx := -1
			for j := i + 1; j < len(tokens); j++ {
				if tokens[j] == ";" {
					endIdx = j
					break
				}
			}

			if endIdx == -1 {
				return fmt.Errorf("missing ; in definition")
			}

			newWord := tokens[i+1]

			if _, err := strconv.Atoi(newWord); err == nil {
				return fmt.Errorf("cannot redefine numbers")
			}

			definition := tokens[i+2 : endIdx]

			expanded := make([]string, 0)
			for _, token := range definition {
				if def, exists := e.words[token]; exists {
					expanded = append(expanded, def...)
				} else if e.isBuiltin(token) {
					// Помечаем встроенные операции
					expanded = append(expanded, "$$"+token)
				} else {
					expanded = append(expanded, token)
				}
			}

			e.words[newWord] = expanded

			i = endIdx + 1
			continue
		}

		// Проверяем префикс встроенной операции
		if strings.HasPrefix(word, "$$") {
			word = strings.TrimPrefix(word, "$$")
			// Выполняем как встроенную операцию (идём к switch)
		} else if def, exists := e.words[word]; exists {
			err := e.processTokens(def)
			if err != nil {
				return err
			}
			i++
			continue
		}

		switch word {
		case "+":
			y, err := e.stack.Pop()
			if err != nil {
				return err
			}
			x, err := e.stack.Pop()
			if err != nil {
				return err
			}
			e.stack.Push(x + y)

		case "-":
			y, err := e.stack.Pop()
			if err != nil {
				return err
			}
			x, err := e.stack.Pop()
			if err != nil {
				return err
			}
			e.stack.Push(x - y)

		case "*":
			y, err := e.stack.Pop()
			if err != nil {
				return err
			}
			x, err := e.stack.Pop()
			if err != nil {
				return err
			}
			e.stack.Push(x * y)

		case "/":
			y, err := e.stack.Pop()
			if err != nil {
				return err
			}
			x, err := e.stack.Pop()
			if err != nil {
				return err
			}
			if y == 0 {
				return fmt.Errorf("division by zero")
			}
			e.stack.Push(x / y)

		case "dup":
			x, err := e.stack.Pop()
			if err != nil {
				return err
			}
			e.stack.Push(x)
			e.stack.Push(x)

		case "over":
			y, err := e.stack.Pop()
			if err != nil {
				return err
			}
			x, err := e.stack.Pop()
			if err != nil {
				return err
			}
			e.stack.Push(x)
			e.stack.Push(y)
			e.stack.Push(x)

		case "drop":
			_, err := e.stack.Pop()
			if err != nil {
				return err
			}

		case "swap":
			y, err := e.stack.Pop()
			if err != nil {
				return err
			}
			x, err := e.stack.Pop()
			if err != nil {
				return err
			}
			e.stack.Push(y)
			e.stack.Push(x)

		default:
			num, err := strconv.Atoi(word)
			if err != nil {
				return fmt.Errorf("unknown word: %s", word)
			}
			e.stack.Push(num)
		}

		i++
	}

	return nil
}
