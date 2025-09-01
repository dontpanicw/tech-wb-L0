package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir(".")) // раздаём файлы из текущей папки
	http.Handle("/", fs)

	addr := ":8000"
	log.Printf("Frontend running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
