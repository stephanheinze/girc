package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Me struct {
	processPattern *regexp.Regexp
	entries        StringList
}

func ME(name string, filename string) *Me {
	me := Me{
		processPattern: regexp.MustCompile(fmt.Sprintf("^:(.*)+\\!+[^ ]+ PRIVMSG ([^ ]+) :.*%s.*$", name)),
	}
	if filename == "" {
		me.addDefaultEntries()
	} else {
		if file, err := os.OpenFile(filename, os.O_RDONLY, 0666); err != nil {
			log.Printf("Could not open 'me'-file %s: %s", filename, err.Error())
			me.addDefaultEntries()
		} else {
			defer file.Close()
			me.entries.Read(file)
		}
	}
	return &me
}

func (this *Me) addDefaultEntries() {
	this.entries.Add("Hey %USER%, stay calm! I'm just a poor bot.")
}

func (this *Me) Match(line string) bool {
	return this.processPattern.MatchString(line)
}

func (this *Me) Process(line string, gossip *Gossip) {
	channel, response := this.process(line)
	gossip.Conn.Cmd("PRIVMSG %s :%s\r\n", channel, response)
}

func (this *Me) process(line string) (channel string, response string) {
	match := this.processPattern.FindStringSubmatch(line)
	channel = match[2]
	_, _, entry := this.entries.Random()
	response = strings.Replace(entry, "%USER%", match[1], -1)
	return
}
