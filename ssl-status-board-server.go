package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
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
		wrapper := new(sslproto.SSL_Micro_WrapperPacket)
		wrapper.Detection = new(sslproto.SSL_Micro_DetectionFrame)
		visionDetectionMutex.Lock()
		removeOldCamDetections()
		for _, r := range latestVisionDetection {
			wrapper.Detection.Balls = append(wrapper.Detection.Balls, micronizeBalls(r.Balls)...)
			wrapper.Detection.RobotsBlue = append(wrapper.Detection.RobotsBlue, micronizeBots(r.RobotsBlue)...)
			wrapper.Detection.RobotsYellow = append(wrapper.Detection.RobotsYellow, micronizeBots(r.RobotsYellow)...)
		}
		visionDetectionMutex.Unlock()
		if latestVisionGeometry != nil && (first || time.Now().Sub(lastTimeGeometrySend) > serverConfig.GeometrySendingInterval) {
			lastTimeGeometrySend = time.Now()
			wrapper.Geometry = micronizeGeometry(latestVisionGeometry)
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

func micronizeBalls(balls []*sslproto.SSL_DetectionBall) (microBalls []*sslproto.SSL_Micro_DetectionBall) {
	microBalls = make([]*sslproto.SSL_Micro_DetectionBall, len(balls))
	for i, b := range balls {
		microBalls[i] = new(sslproto.SSL_Micro_DetectionBall)
		microBalls[i].X = b.X
		microBalls[i].Y = b.Y
	}
	return
}

func micronizeBots(robots []*sslproto.SSL_DetectionRobot) (microRobots []*sslproto.SSL_Micro_DetectionRobot) {
	microRobots = make([]*sslproto.SSL_Micro_DetectionRobot, len(robots))
	for i, r := range robots {
		microRobots[i] = new(sslproto.SSL_Micro_DetectionRobot)
		microRobots[i].RobotId = r.RobotId
		microRobots[i].X = r.X
		microRobots[i].Y = r.Y
		microRobots[i].Orientation = r.Orientation
	}
	return
}

func micronizeGeometry(geometry *sslproto.SSL_GeometryData) (microGeometry *sslproto.SSL_Micro_GeometryData) {
	microGeometry = new(sslproto.SSL_Micro_GeometryData)
	microGeometry.Field = new(sslproto.SSL_Micro_GeometryFieldSize)
	microGeometry.Field.BoundaryWidth = geometry.Field.BoundaryWidth
	microGeometry.Field.FieldLength = geometry.Field.FieldLength
	microGeometry.Field.FieldWidth = geometry.Field.FieldWidth
	microGeometry.Field.GoalDepth = geometry.Field.GoalDepth
	microGeometry.Field.GoalWidth = geometry.Field.GoalWidth
	microGeometry.Field.FieldLines = micronizeLines(geometry.Field.FieldLines)
	microGeometry.Field.FieldArcs = micronizeArcs(geometry.Field.FieldArcs)
	return
}

func micronizeLines(lines []*sslproto.SSL_FieldLineSegment) (microLines []*sslproto.SSL_Micro_FieldLineSegment) {
	microLines = make([]*sslproto.SSL_Micro_FieldLineSegment, len(lines))
	for i, r := range lines {
		microLines[i] = new(sslproto.SSL_Micro_FieldLineSegment)
		microLines[i].P1 = new(sslproto.Micro_Vector2F)
		microLines[i].P1.X = r.P1.X
		microLines[i].P1.Y = r.P1.Y
		microLines[i].P2 = new(sslproto.Micro_Vector2F)
		microLines[i].P2.X = r.P2.X
		microLines[i].P2.Y = r.P2.Y
	}
	return
}

func micronizeArcs(arcs []*sslproto.SSL_FieldCicularArc) (microArcs []*sslproto.SSL_Micro_FieldCicularArc) {
	microArcs = make([]*sslproto.SSL_Micro_FieldCicularArc, len(arcs))
	for i, r := range arcs {
		microArcs[i] = new(sslproto.SSL_Micro_FieldCicularArc)
		microArcs[i].Center = new(sslproto.Micro_Vector2F)
		microArcs[i].Center.X = r.Center.X
		microArcs[i].Center.Y = r.Center.Y
		microArcs[i].Radius = r.Radius
		microArcs[i].A1 = r.A1
		microArcs[i].A2 = r.A2
	}
	return
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
