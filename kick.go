package main

import "regexp"

type Kick struct {
	pattern *regexp.Regexp
}

func KICK() *Kick {
	return &Kick{pattern: regexp.MustCompile("^.* KICK (#[^ ]*).*$")}
}

func (this *Kick) Match(line string) bool {
	return this.pattern.MatchString(line)
}

func (this *Kick) Process(line string, gossip *Gossip) {
	match := this.pattern.FindStringSubmatch(line)
	gossip.Conn.Cmd("PRIVMSG %s :%s\r\n", match[1], "pffft!")
}
