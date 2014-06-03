package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
)

type Cite struct {
	pattern  *regexp.Regexp
	shortCut string
	entries  []string
}

func CITE(shortCut, filename string) *Cite {
	cite := Cite{
		pattern: regexp.MustCompile("^[^ ]+ PRIVMSG ([^ ]+) :!" + shortCut + "$"),
	}
	if file, err := os.Open(filename); err != nil {
		log.Printf("Could not load cite-file %q - reason: %s", filename, err.Error())
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			cite.entries = append(cite.entries, scanner.Text())
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

func (this *Cite) process(line string) (channel string, response string) {
	match := this.pattern.FindStringSubmatch(line)
	channel = match[1]
	max, index, entry := this.Random()
	response = fmt.Sprintf("%s[%d/%d]: %s", this.shortCut, index, max, entry)
	return
}

func (this *Cite) Random() (max int, index int, entry string) {
	max = len(this.entries)
	switch max {
	case 0:
		return 0, 0, "-- no entries --"
	case 1:
		return 1, 1, this.entries[0]
	default:
		index = rand.Intn(max)
		entry = this.entries[index]
		index = index + 1
		return
	}
}
