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

func unmarshallNTPResponse(buffer []byte) (referenceTime NTPTime, receptionTime NTPTime, transmissionTime NTPTime) {
	referenceTime = NTPTimeFromByteArray(buffer[16 : 16+8])
	receptionTime = NTPTimeFromByteArray(buffer[32 : 32+8])
	transmissionTime = NTPTimeFromByteArray(buffer[40 : 40+8])
	return referenceTime, receptionTime, transmissionTime
}

func (client SNTPClient) GetOffset() (time.Duration, error) {
	reqMsg := make([]byte, 48)
	reqMsg[0] = 0x1b
	respMsg := make([]byte, 48)

	clientRequestTransmissionTime := time.Now().UTC()
	_, err := client.Write(reqMsg)
	if err != nil {
		return 0, err
	}

	clientResponseReceptionTime := time.Now().UTC()
	_, err = client.Read(respMsg)
	if err != nil {
		return 0, err
	}

	serverReferenceTime, serverReceptionTime, serverTransmissionTime := unmarshallNTPResponse(respMsg)

	offset := calculateClockOffset(clientRequestTransmissionTime,
		clientResponseReceptionTime,
		serverReceptionTime,
		serverTransmissionTime)

	fmt.Println("Server ref time:", serverReferenceTime.Time().Unix())
	fmt.Println("Local time     :", clientRequestTransmissionTime.Unix())
	fmt.Println("Time difference:", offset)

	return offset, nil
}
