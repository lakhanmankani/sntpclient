package main

import (
	"fmt"
	"github.com/lakhanmankani/sntpclient/client"
	"log"
)

func main() {
	conn, err := client.CreateSNTPConnection("time.google.com")
	if err != nil {
		log.Fatal(err)
	}
	offset, err := conn.GetOffset()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Offset:", offset)

	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
