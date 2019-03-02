package proxy

import (
	"encoding/base64"
	"github.com/RoboCup-SSL/ssl-status-board-server/internal/config"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

func broadcastToProxy(proxyConfig config.ServerProxyConfig, handler func(*websocket.Conn)) error {
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

	handler(conn)

	err = conn.Close()
	if err != nil {
		log.Println("Could not close proxy websocket connection")
	}
	return nil
}

// HandleServerProxy handles proxying of websocket connections
func HandleServerProxy(proxyConfig config.ServerProxyConfig, handler func(*websocket.Conn)) {
	for {
		err := broadcastToProxy(proxyConfig, handler)
		log.Println("Disconnected from proxy ", err)
		if err != nil {
			time.Sleep(proxyConfig.ReconnectInterval)
		}
	}
}
