package memstatsd

import (
	"log"
	"testing"
	"time"
)

type statter struct{}

func (s statter) Timing(bucket string, d time.Duration) {
	log.Println(bucket, d)
}

func (s statter) Gauge(bucket string, value int) {
	log.Println(bucket, value)
}

func TestMemstatsd(t *testing.T) {
	msd := New("memstatsd.", "testing", statter{}, true)
	msd.Run(5 * time.Second)
	time.Sleep(time.Second * 10)

	go func() {
		time.Sleep(time.Minute)
	}()

	time.Sleep(time.Minute)
}
