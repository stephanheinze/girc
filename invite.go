package main

import (
	"fmt"
	"regexp"
)

type Invite struct {
	pattern *regexp.Regexp
}

func INVITE() *Invite {
	invite := Invite{
		pattern: regexp.MustCompile("^:(.*)+\\!+[^ ]+ INVITE [^ ]+ :?(.*)$"),
	}
	return &invite
}

func (this *Invite) Match(line string) bool {
	return this.pattern.MatchString(line)
}

func (this *Invite) Process(line string, gossip *Gossip) {
	channel, name := this.parseLine(line)
	gossip.JoinChannel(channel)
	gossip.SendMessage(channel, fmt.Sprintf("Thanks %s for invitation.", name))
	gossip.AddChannel(channel)
	gossip.PrintChannels()
}

func (this *Invite) parseLine(line string) (channel, name string) {
	match := this.pattern.FindStringSubmatch(line)
	name = match[1]
	channel = match[2]
	return
}
