package main

import (
	"flag"
	"github.com/RoboCup-SSL/ssl-status-board-server/internal/config"
	"github.com/RoboCup-SSL/ssl-status-board-server/internal/proxy"
	"github.com/RoboCup-SSL/ssl-status-board-server/internal/referee"
	"github.com/RoboCup-SSL/ssl-status-board-server/internal/vision"
	"log"
	"net/http"
)

func main() {
	configFile := flag.String("c", "server-config.yaml", "The config file to use")
	flag.Parse()

	serverConfig := config.ReadServerConfig(*configFile)
	log.Println("Server config:", serverConfig)

	refereeBoard := referee.NewBoard(serverConfig.RefereeConnection)
	visionBoard := vision.NewBoard(serverConfig.VisionConnection)

	go refereeBoard.HandleIncomingMessages()
	go visionBoard.HandleIncomingMessages()

	if serverConfig.RefereeConnection.ServerProxy.Enabled {
		go proxy.HandleServerProxy(serverConfig.RefereeConnection.ServerProxy, refereeBoard.SendToWebSocket)
	}
	if serverConfig.VisionConnection.ServerProxy.Enabled {
		go proxy.HandleServerProxy(serverConfig.VisionConnection.ServerProxy, visionBoard.SendToWebSocket)
	}

	http.HandleFunc(serverConfig.RefereeConnection.SubscribePath, refereeBoard.WsHandler)
	http.HandleFunc(serverConfig.VisionConnection.SubscribePath, visionBoard.WsHandler)
	log.Fatal(http.ListenAndServe(serverConfig.ListenAddress, nil))
}
