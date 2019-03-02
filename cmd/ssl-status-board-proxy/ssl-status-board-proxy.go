package main

import (
	"flag"
	"github.com/RoboCup-SSL/ssl-status-board-server/internal/config"
	"github.com/RoboCup-SSL/ssl-status-board-server/internal/proxy"
	"log"
	"net/http"
)

func main() {
	configFile := flag.String("c", "proxy-config.yaml", "The config file to use")
	flag.Parse()

	proxyConfig := config.ReadProxyConfig(*configFile)
	log.Println("Proxy config:", proxyConfig)

	p := proxy.NewProxy(proxyConfig)

	http.HandleFunc(proxyConfig.SubscribePath, p.Serve)
	http.HandleFunc(proxyConfig.PublishPath, p.Receive)
	log.Println("Start listener on", proxyConfig.ListenAddress)
	log.Fatal(http.ListenAndServe(proxyConfig.ListenAddress, nil))
}
