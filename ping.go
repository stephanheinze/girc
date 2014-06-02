package main

import "regexp"

type Ping struct {
	pattern *regexp.Regexp
}

func PING() *Ping {
	return &Ping{pattern: regexp.MustCompile("^PING :([a-zA-Z0-9\\.]+)$")}
}

func (this *Ping) Match(line string) bool {
	return this.pattern.MatchString(line)
}

func (this *Ping) Process(line string, gossip *Gossip) {
	match := this.pattern.FindStringSubmatch(line)
	gossip.Conn.Cmd("PONG :%s\r\n", match[1])
}
