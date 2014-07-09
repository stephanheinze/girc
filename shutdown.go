package main

import (
	"regexp"
)

type Shutdown struct {
	pattern *regexp.Regexp
}

func SHUTDOWN(password string) *Shutdown {
	shutdown := Shutdown{
		pattern: regexp.MustCompile("^:(.*)+\\!+[^ ]+ PRIVMSG .* :shutdown " + password + "$"),
	}
	return &shutdown
}

func (this *Shutdown) Match(line string) bool {
	return this.pattern.MatchString(line)
}

func (this *Shutdown) Process(line string, gossip *Gossip) {
	name := this.parseLine(line)
	gossip.Conn.Cmd("QUIT :"+name+" kicked me - bye bye\r\n", name)
}

func (this *Shutdown) parseLine(line string) (name string) {
	match := this.pattern.FindStringSubmatch(line)
	name = match[1]
	return
}
