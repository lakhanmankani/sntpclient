package client

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

type SNTPTime struct {
	seconds     uint32
	nanoSeconds uint32
}

func SNTPTimeFromByteArray(timeStamp []byte) SNTPTime {
	return SNTPTime{binary.BigEndian.Uint32(timeStamp[0:4]), binary.BigEndian.Uint32(timeStamp[4:])}
}

func (sntpTime *SNTPTime) Time() time.Time {
	time1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	serverTime := time1900.Add(time.Second * time.Duration(sntpTime.seconds))
	serverTime = serverTime.Add(time.Nanosecond * time.Duration(sntpTime.nanoSeconds))

	return serverTime
}

func CreateSocket() {
	udpAddr, _ := net.ResolveUDPAddr("udp", net.JoinHostPort("time.google.com", "123"))
	fmt.Println(udpAddr.String())

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	msg := make([]byte, 48)
	msg[0] = 0x1b
	fmt.Println(msg)
	buffer := make([]byte, 48)

	clientRequestTransmissionTime := time.Now().UTC()
	_, err = conn.Write(msg)
	if err != nil {
		log.Fatal(err)
	}

	clientResponseReceptionTime := time.Now().UTC()
	_, err = conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer)
	fmt.Println(buffer[16 : 16+8])

	serverReferenceTime := SNTPTimeFromByteArray(buffer[16 : 16+8])
	serverReceptionTime := SNTPTimeFromByteArray(buffer[32 : 32+8])
	serverTransmissionTime := SNTPTimeFromByteArray(buffer[40 : 40+8])

	timeDifference := ((serverReceptionTime.Time().Sub(clientRequestTransmissionTime)) +
		(serverTransmissionTime.Time().Sub(clientResponseReceptionTime))) / 2

	fmt.Println("Server ref time:", serverReferenceTime.Time().Unix())
	fmt.Println("Local time     :", clientRequestTransmissionTime.Unix())
	fmt.Println("Time difference:", timeDifference)
}
