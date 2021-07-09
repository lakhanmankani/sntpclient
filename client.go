package main

import (
	"fmt"
	"net"
	"time"
)

type SNTPClient net.UDPConn

type NTPResponse struct {
	referenceTimeStamp NTPTime
	originateTimeStamp NTPTime
	receiveTimeStamp   NTPTime
	transmitTimestamp  NTPTime
}

func CreateSNTPConnection(host string) (*SNTPClient, error) {
	udpAddr, _ := net.ResolveUDPAddr("udp", net.JoinHostPort(host, "123"))
	// fmt.Println(udpAddr.String())

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

func unmarshallNTPResponse(buffer []byte) NTPResponse {
	return NTPResponse{
		referenceTimeStamp: NTPTimeFromByteArray(buffer[16 : 16+8]),
		originateTimeStamp: NTPTimeFromByteArray(buffer[32 : 32+8]),
		receiveTimeStamp:   NTPTimeFromByteArray(buffer[40 : 40+8]),
		transmitTimestamp:  NTPTimeFromByteArray(buffer[48 : 48+8]),
	}
}

func (client SNTPClient) MakeRequest() (receptionTime time.Time, response NTPResponse, err error) {
	reqMsg := make([]byte, 48)
	reqMsg[0] = 0x1b
	respMsg := make([]byte, 56)

	clientRequestTransmissionTime := time.Now().UTC()
	requestNTPTime := NTPTimeFromTime(clientRequestTransmissionTime)
	copy(reqMsg[32:32+8], requestNTPTime.ByteArrayFromNTP())

	fmt.Println(reqMsg)
	// TODO: Write request time to msg

	_, err = client.Write(reqMsg)
	if err != nil {
		return time.Time{}, NTPResponse{}, err
	}

	clientResponseReceptionTime := time.Now().UTC()
	_, err = client.Read(respMsg)
	if err != nil {
		return time.Time{}, NTPResponse{}, err
	}
	fmt.Println(respMsg)

	response = unmarshallNTPResponse(respMsg)


	return clientResponseReceptionTime, response, nil
}

func (client SNTPClient) GetOffset() (time.Duration, error) {
	receptionTime, resp, err := client.MakeRequest()
	if err != nil {
		return 0, err
	}


	offset := calculateClockOffset(resp.originateTimeStamp.Time(),
		receptionTime, // resp.receiveTimeStamp.Time(),
		resp.referenceTimeStamp,
		resp.transmitTimestamp)

	fmt.Println("Server ref time:", resp.referenceTimeStamp.Time().UnixNano())
	fmt.Println("Local time     :", resp.originateTimeStamp.Time().UnixNano())
	fmt.Println("Time difference:", offset)

	return offset, nil
}
