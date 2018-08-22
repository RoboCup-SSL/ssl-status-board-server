package main

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"time"
)

func handleIncomingRefereeMessages() {
	address := serverConfig.RefereeConnection.MulticastAddress
	err, refereeListener := openMulticastUdpConnection(address)
	if err != nil {
		log.Println("Could not connect to ", address)
	}

	lastCommandId := uint32(100000000)
	for {
		data := make([]byte, maxDatagramSize)
		n, _, err := refereeListener.ReadFromUDP(data)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}

		message, err := parseRefereeMessage(data[:n])
		if err != nil {
			log.Print("Could not parse referee message: ", err)
		} else {
			saveRefereeMessageFields(message)

			if *message.CommandCounter != lastCommandId {
				log.Println("Received referee message:", message)
				lastCommandId = *message.CommandCounter
			}
		}
	}
}

func handleIncomingVisionMessages() {
	address := serverConfig.VisionConnection.MulticastAddress
	err, refereeListener := openMulticastUdpConnection(address)
	if err != nil {
		log.Println("Could not connect to ", address)
	}

	for {
		data := make([]byte, maxDatagramSize)
		n, _, err := refereeListener.ReadFromUDP(data)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}

		message, err := parseVisionWrapperPacket(data[:n])
		if err != nil {
			log.Print("Could not parse referee message: ", err)
		} else {
			if message.Detection != nil {
				visionDetectionMutex.Lock()
				camId := int(*message.Detection.CameraId)
				latestVisionDetection[camId] = message.Detection
				visionDetectionReceived[camId] = time.Now()
				visionDetectionMutex.Unlock()
			}
			if message.Geometry != nil {
				latestVisionGeometry = message.Geometry
			}
		}
	}
}

func openMulticastUdpConnection(address string) (err error, listener *net.UDPConn) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	listener, err = net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("could not connect to ", address)
	}
	listener.SetReadBuffer(maxDatagramSize)
	log.Printf("Listening on %s", address)
	return
}

func parseRefereeMessage(data []byte) (message *sslproto.SSL_Referee, err error) {
	message = new(sslproto.SSL_Referee)
	err = proto.Unmarshal(data, message)
	return
}

func parseVisionWrapperPacket(data []byte) (message *sslproto.SSL_WrapperPacket, err error) {
	message = new(sslproto.SSL_WrapperPacket)
	err = proto.Unmarshal(data, message)
	return
}
