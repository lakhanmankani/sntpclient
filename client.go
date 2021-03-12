package client

import (
	"fmt"
	"log"
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
	//defer conn.Close()
	//
	//msg := make([]byte, 48)
	//msg[0] = 0x1b
	//fmt.Println(msg)
	//buffer := make([]byte, 48)
	//
	//clientRequestTransmissionTime := time.Now().UTC()
	//_, err = conn.Write(msg)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//clientResponseReceptionTime := time.Now().UTC()
	//_, err = conn.Read(buffer)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(buffer)
	//fmt.Println(buffer[16 : 16+8])
	//
	//serverReferenceTime := SNTPTimeFromByteArray(buffer[16 : 16+8])
	//serverReceptionTime := SNTPTimeFromByteArray(buffer[32 : 32+8])
	//serverTransmissionTime := SNTPTimeFromByteArray(buffer[40 : 40+8])
	//
	//timeDifference := ((serverReceptionTime.Time().Sub(clientRequestTransmissionTime)) +
	//	(serverTransmissionTime.Time().Sub(clientResponseReceptionTime))) / 2
	//
	//fmt.Println("Server ref time:", serverReferenceTime.Time().Unix())
	//fmt.Println("Local time     :", clientRequestTransmissionTime.Unix())
	//fmt.Println("Time difference:", timeDifference)
	//}

func calculateClockOffset(clientRequestTransmissionTime time.Time,
	clientResponseReceptionTime time.Time,
	serverReceptionTime SNTPTime,
	serverTransmissionTime SNTPTime) time.Duration {
	offset := ((serverReceptionTime.Time().Sub(clientRequestTransmissionTime)) +
		(serverTransmissionTime.Time().Sub(clientResponseReceptionTime))) / 2
	return offset
}

func (client SNTPClient) GetOffset() time.Duration {
	msg := make([]byte, 48)
	msg[0] = 0x1b
	fmt.Println(msg)
	buffer := make([]byte, 48)

	clientRequestTransmissionTime := time.Now().UTC()
	_, err := client.Write(msg)
	if err != nil {
		log.Fatal(err)
	}

	clientResponseReceptionTime := time.Now().UTC()
	_, err = client.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer)
	fmt.Println(buffer[16 : 16+8])

	serverReferenceTime := SNTPTimeFromByteArray(buffer[16 : 16+8])
	serverReceptionTime := SNTPTimeFromByteArray(buffer[32 : 32+8])
	serverTransmissionTime := SNTPTimeFromByteArray(buffer[40 : 40+8])

	offset := calculateClockOffset(clientRequestTransmissionTime,
		clientResponseReceptionTime,
		serverReceptionTime,
		serverTransmissionTime)

	fmt.Println("Server ref time:", serverReferenceTime.Time().Unix())
	fmt.Println("Local time     :", clientRequestTransmissionTime.Unix())
	fmt.Println("Time difference:", offset)

	return offset
}
