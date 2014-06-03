package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Cite struct {
	pattern  *regexp.Regexp
	shortCut string
	entries  StringList
}

func CITE(shortCut, filename string) *Cite {
	cite := Cite{
		pattern:  regexp.MustCompile("^[^ ]+ PRIVMSG ([^ ]+) :!" + shortCut + " *([^ ]*)$"),
		shortCut: shortCut,
	}
	if file, err := os.Open(filename); err != nil {
		log.Printf("Could not load cite-file %q - reason: %s", filename, err.Error())
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			cite.entries.Add(scanner.Text())
		}
	}
	return &cite
}

func (this *Cite) Match(line string) bool {
	return this.pattern.MatchString(line)
}

func (this *Cite) Process(line string, gossip *Gossip) {
	channel, response := this.process(line)
	gossip.Conn.Cmd("PRIVMSG %s :%s\r\n", channel, response)
}

func (this *Cite) process(line string) (channel, response string) {
	match := this.pattern.FindStringSubmatch(line)
	channel = match[1]
	filter := match[2]
	var (
		max, index int
		entry      string
	)
	if filter != "" {
		max, index, entry = this.entries.FilteredRandom(filter)
	} else {
		max, index, entry = this.entries.Random()
	}
	response = fmt.Sprintf("%s[%d/%d]: %s", this.shortCut, index, max, entry)
	return
}
