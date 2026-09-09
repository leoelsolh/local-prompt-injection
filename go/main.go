package main

import (
	"fmt"
	"log"
)

func main() {
	content, err := runAgent("qwen2.5:1.5b", "What does the page at http://127.0.0.1:5000 say?")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(content)
}
