package client

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

func NTPTimeFromTime(t time.Time) NTPTime {
	time1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	timeDiff := t.Sub(time1900)
	return NTPTime{uint32(timeDiff.Seconds()), uint32(timeDiff.Nanoseconds())}
}

func (ntpTime *NTPTime) Time() time.Time {
	time1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	serverTime := time1900.Add(time.Second * time.Duration(ntpTime.seconds))
	serverTime = serverTime.Add(time.Nanosecond * time.Duration(ntpTime.nanoSeconds))

	return serverTime
}

func (ntpTime *NTPTime) ByteArrayFromNTP() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b[0:4], ntpTime.seconds)
	binary.BigEndian.PutUint32(b[4:8], ntpTime.nanoSeconds)
	return b
}
