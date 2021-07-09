package main

import (
	"fmt"
	"log"
)

func main() {
	client, err := CreateSNTPConnection("time.google.com")
	if err != nil {
		log.Fatal(err)
	}
	offset, err := client.GetOffset()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Offset:", offset)

	err = client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
