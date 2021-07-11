package client

import (
	"encoding/binary"
	"math"
	"time"
)

type NTPTime struct {
	Seconds  uint32
	Fraction uint32
}

func NTPTimeFromByteArray(timeStamp []byte) NTPTime {
	return NTPTime{binary.BigEndian.Uint32(timeStamp[0:4]), binary.BigEndian.Uint32(timeStamp[4:8])}
}

func NTPTimeFromTime(t time.Time) NTPTime {
	time1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	timeDiff := t.Sub(time1900)
	seconds := timeDiff.Seconds()
	fraction := (seconds - math.Trunc(seconds)) * math.MaxUint32
	return NTPTime{uint32(seconds), uint32(fraction)}
}

func (ntpTime *NTPTime) Time() time.Time {
	time1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	serverTime := time1900.Add(time.Second * time.Duration(ntpTime.Seconds))
	serverTime = serverTime.Add(time.Second * time.Duration(ntpTime.Fraction) / math.MaxUint32)

	return serverTime
}

func (ntpTime *NTPTime) ByteArrayFromNTP() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b[0:4], ntpTime.Seconds)
	binary.BigEndian.PutUint32(b[4:8], ntpTime.Fraction)
	return b
}
