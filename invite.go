package main

import (
	"regexp"
)

type Invite struct {
	pattern *regexp.Regexp
}

func INVITE() *Invite {
	invite := Invite{
		pattern: regexp.MustCompile("^:(.*)+\\!+[^ ]+ INVITE [^ ]+ :(.*)$"),
	}
	return &invite
}

func (this *Invite) Match(line string) bool {
	return this.pattern.MatchString(line)
}

func (this *Invite) Process(line string, gossip *Gossip) {
	channel, name := this.parseLine(line)
	gossip.Conn.Cmd("JOIN %s\r\n", channel)
	gossip.Conn.Cmd("PRIVMSG %s :Thanks %s for invitation.\r\n", channel, name)
}

func (this *Invite) parseLine(line string) (channel, name string) {
	match := this.pattern.FindStringSubmatch(line)
	name = match[1]
	channel = match[2]
	return
}
