package main

import (
	"flag"
	"ledz/spi"
	"log"
	"net"
	"time"
)

var (
	speed     = flag.Int64("speed", 2000000, "spi speed")
	dev       = flag.String("device", "/dev/spidev0.0", "spi device")
	frameRate = flag.Int64("rate", 60, "framerate")
	size      = flag.Int64("size", 150, "size of frame")
	listen    = flag.String("listen", ":1337", "bind address")
	buffer    = []byte{}
)

func renderer(device *spi.Device) {
	ticker := time.NewTicker(time.Second / time.Duration(*frameRate))
	for range ticker.C {
		if _, err := device.Tx(buffer); err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {
	var con *spi.Device
	var err error
	flag.Parse()
	buffer = make([]byte, int(*size))

	for con == nil {
		con, err = spi.Open(*dev, &spi.Config{
			Mode:          0,
			Bits:          8,
			Speed:         uint32(*speed),
			DelayUsec:     500,
			CSChange:      false,
			TXNBits:       0,
			RXNBits:       0,
			WordDelayUsec: 0,
		})

		if err != nil {
			log.Println("Open:", err)
			time.Sleep(time.Second)
		}
	}
	go renderer(con)
	sock, err := net.ListenPacket("udp", *listen)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		_, _, err := sock.ReadFrom(buffer)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
