package vision

import (
	"fmt"
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
	"github.com/RoboCup-SSL/ssl-status-board-server/internal/board"
	"github.com/RoboCup-SSL/ssl-status-board-server/internal/config"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

// Board contains the state of this vision board
type Board struct {
	cfg               config.VisionConnection
	detection         map[int]*sslproto.SSL_DetectionFrame
	detectionReceived map[int]time.Time
	detectionMutex    sync.Mutex
	geometry          *sslproto.SSL_GeometryData
	geometryLastSend  time.Time
}

// NewBoard creates a new vision board
func NewBoard(cfg config.VisionConnection) (b Board) {
	b.cfg = cfg
	b.detection = map[int]*sslproto.SSL_DetectionFrame{}
	b.detectionReceived = map[int]time.Time{}
	b.geometryLastSend = time.Now()
	return
}

// HandleIncomingMessages listens for new messages and stores the latest ones
func (b *Board) HandleIncomingMessages() {
	board.HandleIncomingMessages(b.cfg.ConnectionConfig, b.handlingMessage)
}

func (b *Board) handlingMessage(data []byte) {
	message := new(sslproto.SSL_WrapperPacket)
	err := proto.Unmarshal(data, message)
	if err != nil {
		log.Print("Could not parse referee message: ", err)
	} else {
		if message.Detection != nil {
			b.detectionMutex.Lock()
			camId := int(*message.Detection.CameraId)
			b.detection[camId] = message.Detection
			b.detectionReceived[camId] = time.Now()
			b.detectionMutex.Unlock()
		}
		if message.Geometry != nil {
			b.geometry = message.Geometry
		}
	}
}

// SendToWebSocket sends latest data to the given websocket
func (b *Board) SendToWebSocket(conn *websocket.Conn) {
	first := true
	for {
		wrapper := new(sslproto.SSL_Micro_WrapperPacket)
		wrapper.Detection = new(sslproto.SSL_Micro_DetectionFrame)
		b.detectionMutex.Lock()
		b.removeOldCamDetections()
		for _, r := range b.detection {
			wrapper.Detection.Balls = append(wrapper.Detection.Balls, micronizeBalls(r.Balls)...)
			wrapper.Detection.RobotsBlue = append(wrapper.Detection.RobotsBlue, micronizeBots(r.RobotsBlue)...)
			wrapper.Detection.RobotsYellow = append(wrapper.Detection.RobotsYellow, micronizeBots(r.RobotsYellow)...)
		}
		b.detectionMutex.Unlock()
		if b.geometry != nil && (first || time.Now().Sub(b.geometryLastSend) > b.cfg.GeometrySendingInterval) {
			b.geometryLastSend = time.Now()
			wrapper.Geometry = micronizeGeometry(b.geometry)
			first = false
		}

		data, err := proto.Marshal(wrapper)
		if err != nil {
			fmt.Println("Marshal error:", err)
		} else if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Println("Could not write to vision websocket: ", err)
			return
		}

		time.Sleep(b.cfg.SendingInterval)
	}
}

func (b *Board) removeOldCamDetections() {
	for camId, r := range b.detectionReceived {
		if time.Now().Sub(r) > time.Second {
			delete(b.detectionReceived, camId)
		}
	}
}

// WsHandler handles vision websocket connections
func (b *Board) WsHandler(w http.ResponseWriter, r *http.Request) {
	board.WsHandler(w, r, b.SendToWebSocket)
}
