//go:build !solution

package ciletters

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed letter.tmpl
var letterTemplate string

func MakeLetter(n *Notification) (string, error) {
	funcMap := template.FuncMap{
		"commitHash": func(hash string) string {
			if len(hash) >= 8 {
				return hash[:8]
			}
			return hash
		},
		"lastNLines": func(log string, n int) string {
			lines := strings.Split(log, "\n")
			if len(lines) <= n {
				return log
			}
			return strings.Join(lines[len(lines)-n:], "\n")
		},
		"indent12": func(s string) string {
			lines := strings.Split(s, "\n")
			for i, line := range lines {
				if line != "" {
					lines[i] = "            " + line
				}
			}
			return strings.Join(lines, "\n")
		},
	}

	tmpl, err := template.New("letter").Funcs(funcMap).Parse(letterTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, n)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
