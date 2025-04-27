package main

import (
	"log"
	"os"
)

func main() {
	if err := os.Mkdir("uploads", 0600); os.IsNotExist(err) {
		log.Println(err)
	}

	if err := os.Mkdir("templates", 0700); os.IsNotExist(err) {
		log.Println(err)
	}
}
