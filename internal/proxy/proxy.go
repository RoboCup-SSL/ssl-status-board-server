package proxy

import (
	"github.com/RoboCup-SSL/ssl-status-board-server/internal/config"
)

type Proxy struct {
	cfg                      config.ProxyConfig
	messageChannels          []chan WsMessage
	messageProviderConnected bool
	credentials              []string
}

func NewProxy(cfg config.ProxyConfig) (p Proxy) {
	p.cfg = cfg
	p.messageProviderConnected = false
	p.credentials = []string{}
	p.loadCredentials()
	return
}
