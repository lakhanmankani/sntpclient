package client

import (
	"net"
	"time"
)

type SNTPClient net.UDPConn

type NTPResponse struct {
	OriginateTimeStamp NTPTime // T1
	ReceiveTimeStamp   NTPTime // T2
	TransmitTimeStamp  NTPTime // T3
}

func CreateSNTPConnection(host string) (*SNTPClient, error) {
	udpAddr, _ := net.ResolveUDPAddr("udp", net.JoinHostPort(host, "123"))

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	return (*SNTPClient)(conn), err
}

func CalculateClockOffset(resp NTPResponse,
	clientResponseReceptionTime time.Time) time.Duration {
	offset := ((resp.ReceiveTimeStamp.Time().Sub(resp.OriginateTimeStamp.Time())) +
		(resp.TransmitTimeStamp.Time().Sub(clientResponseReceptionTime))) / 2
	return offset
}

func unmarshallNTPResponse(buffer []byte) NTPResponse {
	return NTPResponse{
		OriginateTimeStamp: NTPTimeFromByteArray(buffer[24 : 24+8]),
		ReceiveTimeStamp:   NTPTimeFromByteArray(buffer[32 : 32+8]),
		TransmitTimeStamp:  NTPTimeFromByteArray(buffer[40 : 40+8]),
	}
}

func (client SNTPClient) MakeRequest() (receptionTime time.Time, response NTPResponse, err error) {
	reqMsg := make([]byte, 56)
	reqMsg[0] = 0x1b
	respMsg := make([]byte, 56)

	clientRequestTransmissionTime := time.Now().UTC()
	requestNTPTime := NTPTimeFromTime(clientRequestTransmissionTime)
	copy(reqMsg[40:40+8], requestNTPTime.ByteArrayFromNTP())

	_, err = client.Write(reqMsg)
	if err != nil {
		return time.Time{}, NTPResponse{}, err
	}

	_, err = client.Read(respMsg)
	if err != nil {
		return time.Time{}, NTPResponse{}, err
	}
	clientResponseReceptionTime := time.Now().UTC() // T4

	response = unmarshallNTPResponse(respMsg)

	return clientResponseReceptionTime, response, nil
}

func (client SNTPClient) GetOffset() (time.Duration, error) {
	receptionTime, resp, err := client.MakeRequest()
	if err != nil {
		return 0, err
	}

	offset := CalculateClockOffset(resp, receptionTime)

	return offset, nil
}
