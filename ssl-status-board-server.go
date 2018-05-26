package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/RoboCup-SSL/ssl-go-tools/sslproto"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const maxDatagramSize = 8192

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

var serverConfig ServerConfig
var latestVisionDetection = map[int]*sslproto.SSL_DetectionFrame{}
var visionDetectionReceived = map[int]time.Time{}
var latestVisionGeometry *sslproto.SSL_GeometryData
var lastTimeGeometrySend = time.Now()
var visionDetectionMutex = &sync.Mutex{}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	defer log.Println("Client disconnected")

	log.Println("Client connected")

	sendRefereeDataToWebSocket(conn)
}

func visionHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	defer log.Println("Client disconnected")

	log.Println("Client connected")

	sendVisionDataToWebSocket(conn)
}

func sendRefereeDataToWebSocket(conn *websocket.Conn) {
	for {
		b, err := json.Marshal(referee)
		if err != nil {
			fmt.Println("Marshal error:", err)
		}
		if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
			log.Println(err)
			return
		}

		time.Sleep(serverConfig.RefereeConnection.SendingInterval)
	}
}

func sendVisionDataToWebSocket(conn *websocket.Conn) {
	first := true
	for {
		wrapper := new(sslproto.SSL_WrapperPacket)
		wrapper.Detection = new(sslproto.SSL_DetectionFrame)
		wrapper.Detection.CameraId = new(uint32)
		wrapper.Detection.FrameNumber = new(uint32)
		wrapper.Detection.TCapture = new(float64)
		wrapper.Detection.TSent = new(float64)
		visionDetectionMutex.Lock()
		removeOldCamDetections()
		for _, r := range latestVisionDetection {
			*wrapper.Detection.TCapture = math.Max(*wrapper.Detection.TCapture, *r.TCapture)
			*wrapper.Detection.TSent = math.Max(*wrapper.Detection.TSent, *r.TSent)
			*wrapper.Detection.FrameNumber = uint32(math.Max(float64(*wrapper.Detection.FrameNumber), float64(*r.FrameNumber)))
			wrapper.Detection.Balls = append(wrapper.Detection.Balls, r.Balls...)
			wrapper.Detection.RobotsBlue = append(wrapper.Detection.RobotsBlue, r.RobotsBlue...)
			wrapper.Detection.RobotsYellow = append(wrapper.Detection.RobotsYellow, r.RobotsYellow...)
		}
		visionDetectionMutex.Unlock()
		if first || (latestVisionGeometry != nil && time.Now().Sub(lastTimeGeometrySend) > serverConfig.GeometrySendingInterval) {
			lastTimeGeometrySend = time.Now()
			wrapper.Geometry = latestVisionGeometry
			first = false
		}

		b, err := proto.Marshal(wrapper)
		if err != nil {
			fmt.Println("Marshal error:", err)
		} else if err := conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
			log.Println(err)
			return
		}

		time.Sleep(serverConfig.VisionConnection.SendingInterval)
	}
}

func removeOldCamDetections() {
	for camId, r := range visionDetectionReceived {
		if time.Now().Sub(r) > time.Second {
			delete(visionDetectionReceived, camId)
		}
	}
}

func broadcastToProxy(proxyConfig ServerProxyConfig, handler func(*websocket.Conn)) error {
	u := url.URL{Scheme: proxyConfig.Scheme, Host: proxyConfig.Address, Path: proxyConfig.Path}
	log.Printf("connecting to %s", u.String())

	auth := []byte(proxyConfig.User + ":" + proxyConfig.Password)
	authBase64 := base64.StdEncoding.EncodeToString(auth)

	requestHeader := http.Header{}
	requestHeader.Set("Authorization", "Basic "+authBase64)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), requestHeader)
	if err != nil {
		return err
	}
	defer conn.Close()

	handler(conn)
	return nil
}

func handleServerProxy(proxyConfig ServerProxyConfig, handler func(*websocket.Conn)) {
	for {
		err := broadcastToProxy(proxyConfig, handler)
		log.Println("Disconnected from proxy ", err)
		if err != nil {
			time.Sleep(proxyConfig.ReconnectInterval)
		}
	}
}

func main() {

	configFile := flag.String("c", "server-config.yaml", "The config file to use")
	flag.Parse()

	serverConfig = ReadServerConfig(*configFile)
	log.Println("Server config:", serverConfig)

	go handleIncomingRefereeMessages()
	go handleIncomingVisionMessages()

	if serverConfig.RefereeConnection.ServerProxy.Enabled {
		go handleServerProxy(serverConfig.RefereeConnection.ServerProxy, sendRefereeDataToWebSocket)
	}
	if serverConfig.VisionConnection.ServerProxy.Enabled {
		go handleServerProxy(serverConfig.VisionConnection.ServerProxy, sendVisionDataToWebSocket)
	}

	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc(serverConfig.RefereeConnection.SubscribePath, statusHandler)
	http.HandleFunc(serverConfig.VisionConnection.SubscribePath, visionHandler)
	log.Fatal(http.ListenAndServe(serverConfig.ListenAddress, nil))
}
