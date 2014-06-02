package main

import "regexp"

type Join struct {
	pattern *regexp.Regexp
}

func JOIN() *Join {
	return &Join{pattern: regexp.MustCompile(":End of /MOTD command\\.$")}
}

func (this *Join) Match(line string) bool {
	return this.pattern.MatchString(line)
}

func (this *Join) Process(_ string, gossip *Gossip) {
	gossip.Conn.Cmd("JOIN %s\r\n", gossip.Channel)
}
