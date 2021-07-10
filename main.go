package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lakhanmankani/sntpclient/client"
)

func main() {
	var server string

	if len(os.Args) > 1 {
		server = os.Args[1]
	} else {
		server = "time.google.com"
	}

	conn, err := client.CreateSNTPConnection(server)
	if err != nil {
		log.Fatal(err)
	}
	offset, err := conn.GetOffset()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(server, conn.RemoteAddr())
	fmt.Println("Offset:", offset)

	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
