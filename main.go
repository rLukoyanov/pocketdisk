package main

import (
	"log"
	"os"
)

func main() {
	if err := os.Mkdir("upload", 0600); err != nil {
		log.Println(err)
	}

	if err := os.Mkdir("templates", 0700); err != nil {
		log.Println(err)
	}
}
