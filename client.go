package main

import (
	"fmt"
	"net"
	"time"
)

type SNTPClient net.UDPConn

func CreateSNTPConnection(host string) (*SNTPClient, error) {
	udpAddr, _ := net.ResolveUDPAddr("udp", net.JoinHostPort(host, "123"))
	fmt.Println(udpAddr.String())

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		// log.Fatalf("Failed to dial: %v", err)
		return nil, err
	}
	return (*SNTPClient)(conn), err
}

func calculateClockOffset(clientRequestTransmissionTime time.Time,
	clientResponseReceptionTime time.Time,
	serverReceptionTime NTPTime,
	serverTransmissionTime NTPTime) time.Duration {
	offset := ((serverReceptionTime.Time().Sub(clientRequestTransmissionTime)) +
		(serverTransmissionTime.Time().Sub(clientResponseReceptionTime))) / 2
	return offset
}

func (client SNTPClient) GetOffset() (time.Duration, error) {
	msg := make([]byte, 48)
	msg[0] = 0x1b
	fmt.Println(msg)
	buffer := make([]byte, 48)

	clientRequestTransmissionTime := time.Now().UTC()
	_, err := client.Write(msg)
	if err != nil {
		return 0, err
	}

	clientResponseReceptionTime := time.Now().UTC()
	_, err = client.Read(buffer)
	if err != nil {
		return 0, err
	}
	fmt.Println(buffer)
	fmt.Println(buffer[16 : 16+8])

	serverReferenceTime := NTPTimeFromByteArray(buffer[16 : 16+8])
	serverReceptionTime := NTPTimeFromByteArray(buffer[32 : 32+8])
	serverTransmissionTime := NTPTimeFromByteArray(buffer[40 : 40+8])

	offset := calculateClockOffset(clientRequestTransmissionTime,
		clientResponseReceptionTime,
		serverReceptionTime,
		serverTransmissionTime)

	fmt.Println("Server ref time:", serverReferenceTime.Time().Unix())
	fmt.Println("Local time     :", clientRequestTransmissionTime.Unix())
	fmt.Println("Time difference:", offset)

	return offset, nil
}
