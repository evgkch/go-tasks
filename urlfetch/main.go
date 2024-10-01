package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: urlfetch <url1> <url2> ...")
		os.Exit(1)
	}

	for _, rawURL := range args {
		// Проверяем валидность URL перед запросом
		if _, err := url.ParseRequestURI(rawURL); err != nil {
			fmt.Fprintf(os.Stderr, "fetch: invalid URL %q: %v\n", rawURL, err)
			os.Exit(1)
		}

		// Выполняем HTTP GET запрос
		resp, err := http.Get(rawURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		// Копируем содержимое в stdout
		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", rawURL, err)
			resp.Body.Close()
			os.Exit(1)
		}

		resp.Body.Close()

		// Добавляем перенос строки между разными URL, если нужно
		if len(args) > 1 {
			fmt.Println()
		}
	}
}
