package proxy

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func (p *Proxy) Receive(w http.ResponseWriter, r *http.Request) {

	user, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("You have to authenticate yourself."))
		logFailedToAnswer(err)
		log.Println("Message provider tried to connect without credentials")
		return
	}
	if !p.validCredentials(user, password) {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Your credentials are invalid."))
		logFailedToAnswer(err)
		log.Printf("Message provider tried to connect with wrong credentials: %s:%s\n", user, password)
		return
	}

	if p.messageProviderConnected {
		w.WriteHeader(http.StatusConflict)
		_, err := w.Write([]byte("There is already a message provider connected!"))
		logFailedToAnswer(err)
		log.Println("Another message provider tried to connect")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Message provider connected")

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("Could not read from message provider: ", err)
			break
		}
		for _, c := range p.messageChannels {
			c <- WsMessage{messageType, data}
		}
	}

	p.disconnectMessageProvider(conn)
}

func logFailedToAnswer(err error) {
	if err != nil {
		log.Println("Could not respond to message provider request")
	}
}

func (p *Proxy) disconnectMessageProvider(conn *websocket.Conn) {
	err := conn.Close()
	if err != nil {
		log.Println("Could not close connection to message provider: ", err)
	}
	p.messageProviderConnected = false
	log.Println("Status provider disconnected")
}
