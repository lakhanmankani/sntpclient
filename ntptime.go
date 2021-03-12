package main

import (
	"encoding/binary"
	"time"
)

type NTPTime struct {
	seconds     uint32
	nanoSeconds uint32
}

func NTPTimeFromByteArray(timeStamp []byte) NTPTime {
	return NTPTime{binary.BigEndian.Uint32(timeStamp[0:4]), binary.BigEndian.Uint32(timeStamp[4:])}
}

func (sntpTime *NTPTime) Time() time.Time {
	time1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	serverTime := time1900.Add(time.Second * time.Duration(sntpTime.seconds))
	serverTime = serverTime.Add(time.Nanosecond * time.Duration(sntpTime.nanoSeconds))

	return serverTime
}
