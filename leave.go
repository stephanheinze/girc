package main

import (
	"fmt"
	"regexp"
)

type Leave struct {
	pattern *regexp.Regexp
}

func LEAVE() *Leave {
	leave := Leave{
		pattern: regexp.MustCompile("^:(.*)+\\!+[^ ]+ PRIVMSG ([^ ]+) :#kick.*$"),
	}
	return &leave
}

func (this *Leave) Match(line string) bool {
	return this.pattern.MatchString(line)
}

func (this *Leave) Process(line string, gossip *Gossip) {
	channel, name := this.parseLine(line)
	gossip.LeaveChannel(channel, fmt.Sprintf("reason: %s wants me to leave", name))
}

func (this *Leave) parseLine(line string) (channel, name string) {
	match := this.pattern.FindStringSubmatch(line)
	name = match[1]
	channel = match[2]
	return
}
