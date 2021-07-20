# sntpclient
Very simple SNTP client.

## Installation
```bash
$ go get -u github.com/lakhanmankani/sntpclient
```

## Usage
```bash
$ sntpclient

time.google.com 216.239.35.0:123
Offset: 5.437645ms
```

```bash
$ sntpclient uk.pool.ntp.org

uk.pool.ntp.org 195.171.43.12:123
Offset: -20.245718ms
```

## Example API usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/lakhanmankani/sntpclient/client"
)

func main() {
	conn, err := client.CreateSNTPConnection("time.google.com")
	if err != nil {
		log.Fatal(err)
	}

	receptionTime, resp, err := conn.MakeRequest()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Request sent time:    ", resp.OriginateTimeStamp.Time())
	fmt.Println("Server receive time:  ", resp.ReceiveTimeStamp.Time())
	fmt.Println("Server transmit time: ", resp.TransmitTimeStamp.Time())
	fmt.Println("Response receive time:", receptionTime)

	offset := client.CalculateClockOffset(resp, receptionTime)
	fmt.Println("Offset:", offset)

	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
```
