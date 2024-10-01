//go:build !solution

package speller

import "strings"

var ones = map[int64]string{
	0:  "zero",
	1:  "one",
	2:  "two",
	3:  "three",
	4:  "four",
	5:  "five",
	6:  "six",
	7:  "seven",
	8:  "eight",
	9:  "nine",
	10: "ten",
	11: "eleven",
	12: "twelve",
	13: "thirteen",
	14: "fourteen",
	15: "fifteen",
	16: "sixteen",
	17: "seventeen",
	18: "eighteen",
	19: "nineteen",
}

// Десятки 20-90
var tens = map[int64]string{
	20: "twenty",
	30: "thirty",
	40: "forty",
	50: "fifty",
	60: "sixty",
	70: "seventy",
	80: "eighty",
	90: "ninety",
}

// Крупные разряды
var scales = []struct {
	value int64
	name  string
}{
	{1000000000, "billion"},
	{1000000, "million"},
	{1000, "thousand"},
}

func Spell(n int64) string {
	if n == 0 {
		return "zero"
	}

	if n < 0 {
		return "minus " + Spell(-n)
	}

	var parts []string

	// Обрабатываем крупные разряды
	for _, scale := range scales {
		if n >= scale.value {
			quotient := n / scale.value
			parts = append(parts, spellUnder1000(quotient)+" "+scale.name)
			n = n % scale.value
		}
	}

	// Обрабатываем остаток < 1000
	if n > 0 {
		parts = append(parts, spellUnder1000(n))
	}

	return strings.Join(parts, " ")
}

// spellUnder1000 преобразует число от 1 до 999 в слова
func spellUnder1000(n int64) string {
	if n == 0 {
		return ""
	}

	var parts []string

	// Сотни
	if n >= 100 {
		hundreds := n / 100
		parts = append(parts, ones[hundreds]+" hundred")
		n = n % 100
	}

	// Десятки и единицы
	if n > 0 {
		if n < 20 {
			// 1-19: прямой поиск
			parts = append(parts, ones[n])
		} else {
			// 20-99
			tensDigit := (n / 10) * 10
			onesDigit := n % 10
			if onesDigit == 0 {
				parts = append(parts, tens[tensDigit])
			} else {
				parts = append(parts, tens[tensDigit]+"-"+ones[onesDigit])
			}
		}
	}

	return strings.Join(parts, " ")
}
