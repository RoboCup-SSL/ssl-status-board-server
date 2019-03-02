package proxy

func (p *Proxy) loadCredentials() {
	for _, a := range p.cfg.AuthCredentials {
		p.credentials = append(p.credentials, a.Username+":"+a.Password)
	}
}

func (p *Proxy) validCredentials(user string, password string) bool {
	credential := user + ":" + password
	for _, c := range p.credentials {
		if c == credential {
			return true
		}
	}
	return false
}
