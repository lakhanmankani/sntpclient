package main

import (
	"fmt"
	"net"
	"time"
)

type SNTPClient net.UDPConn

type NTPResponse struct {
	originateTimeStamp NTPTime
	receiveTimeStamp   NTPTime
	transmitTimestamp  NTPTime
}

func CreateSNTPConnection(host string) (*SNTPClient, error) {
	udpAddr, _ := net.ResolveUDPAddr("udp", net.JoinHostPort(host, "123"))

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	return (*SNTPClient)(conn), err
}

func calculateClockOffset(resp NTPResponse,
	clientResponseReceptionTime time.Time) time.Duration {
	fmt.Println("Up:", resp.receiveTimeStamp.Time().Sub(resp.originateTimeStamp.Time()))
	fmt.Println("Down:", resp.transmitTimestamp.Time().Sub(clientResponseReceptionTime))
	offset := ((resp.receiveTimeStamp.Time().Sub(resp.originateTimeStamp.Time())) +
		(resp.transmitTimestamp.Time().Sub(clientResponseReceptionTime))) / 2
	return offset
}

func unmarshallNTPResponse(buffer []byte) NTPResponse {
	return NTPResponse{
		originateTimeStamp: NTPTimeFromByteArray(buffer[16 : 16+8]),
		receiveTimeStamp:   NTPTimeFromByteArray(buffer[32 : 32+8]),
		transmitTimestamp:  NTPTimeFromByteArray(buffer[40 : 40+8]),
	}
}

func (client SNTPClient) MakeRequest() (receptionTime time.Time, response NTPResponse, err error) {
	reqMsg := make([]byte, 56)
	reqMsg[0] = 0x1b
	respMsg := make([]byte, 56)

	clientRequestTransmissionTime := time.Now().UTC()
	requestNTPTime := NTPTimeFromTime(clientRequestTransmissionTime)
	copy(reqMsg[16:16+8], requestNTPTime.ByteArrayFromNTP())

	_, err = client.Write(reqMsg)
	if err != nil {
		return time.Time{}, NTPResponse{}, err
	}

	clientResponseReceptionTime := time.Now().UTC()
	_, err = client.Read(respMsg)
	if err != nil {
		return time.Time{}, NTPResponse{}, err
	}

	response = unmarshallNTPResponse(respMsg)

	return clientResponseReceptionTime, response, nil
}

func (client SNTPClient) GetOffset() (time.Duration, error) {
	receptionTime, resp, err := client.MakeRequest()
	if err != nil {
		return 0, err
	}

	offset := calculateClockOffset(resp, receptionTime)

	fmt.Println("Server rec time:", resp.receiveTimeStamp.Time())
	fmt.Println("Local time     :", resp.originateTimeStamp.Time())
	fmt.Println("Time difference:", offset)

	return offset, nil
}
