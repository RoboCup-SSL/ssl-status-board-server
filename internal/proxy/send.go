package proxy

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func (p *Proxy) Serve(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go p.sendMessages(conn)
}

func (p *Proxy) sendMessages(conn *websocket.Conn) {
	c := make(chan WsMessage, 10)
	p.messageChannels = append(p.messageChannels, c)
	log.Printf("Client connected, now %d clients.\n", len(p.messageChannels))

	for {
		wsMsg := <-c

		if err := conn.WriteMessage(wsMsg.messageType, wsMsg.data); err != nil {
			log.Println("Could not write to message consumer", err)
			p.messageChannels = remove(p.messageChannels, c)
			break
		}
	}

	err := conn.Close()
	if err != nil {
		log.Println("Could not close connection to message consumer")
	}
}

func remove(in []chan WsMessage, rem chan WsMessage) (out []chan WsMessage) {
	out = []chan WsMessage{}
	for _, c := range in {
		if rem != c {
			out = append(out, c)
		}
	}
	return
}
