package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

type WsMessage struct {
	messageType int
	data        []byte
}

var messageChannels []chan WsMessage
var messageProviderConnected = false
var credentials []string
var proxyConfig ProxyConfig

func receiveHandler(w http.ResponseWriter, r *http.Request) {

	user, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You have to authenticate yourself."))
		log.Println("Message provider tried to connect without credentials")
		return
	}
	if !validCredentials(user, password) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Your credentials are invalid."))
		log.Println("Message provider tried to connect with wrong credentials:", user, password)
		return
	}

	if messageProviderConnected {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("There is already a message provider connected!"))
		log.Println("Another message provider tried to connect")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer disconnectMessageProvider(conn)

	log.Println("Message provider connected")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		for _, c := range messageChannels {
			c <- WsMessage{messageType, p}
		}
	}
}

func validCredentials(user string, password string) bool {
	credential := user + ":" + password
	for _, c := range credentials {
		if c == credential {
			return true
		}
	}
	return false
}

func disconnectMessageProvider(conn *websocket.Conn) {
	log.Println("Status provider disconnected")
	messageProviderConnected = false
	conn.Close()
}

func sendMessages(conn *websocket.Conn) {
	c := make(chan WsMessage, 10)
	messageChannels = append(messageChannels, c)
	log.Printf("Client connected, now %d clients.\n", len(messageChannels))

	for {
		wsMsg := <-c

		if err := conn.WriteMessage(wsMsg.messageType, wsMsg.data); err != nil {
			log.Println(err)
			messageChannels = remove(messageChannels, c)
			conn.Close()
			return
		}
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

func serveHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go sendMessages(conn)
}

func loadCredentials() {
	for _, a := range proxyConfig.AuthCredentials {
		credentials = append(credentials, a.Username+":"+a.Password)
	}
}

func main() {

	configFile := flag.String("c", "proxy-config.yaml", "The config file to use")
	flag.Parse()

	proxyConfig = ReadProxyConfig(*configFile)
	log.Println("Proxy config:", proxyConfig)

	loadCredentials()

	http.HandleFunc(proxyConfig.SubscribePath, serveHandler)
	http.HandleFunc(proxyConfig.PublishPath, receiveHandler)
	log.Println("Start listener on", proxyConfig.ListenAddress)
	log.Fatal(http.ListenAndServe(proxyConfig.ListenAddress, nil))
}
