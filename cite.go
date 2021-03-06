package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Cite struct {
	commandPattern *regexp.Regexp
	indexPattern   *regexp.Regexp
	processPattern *regexp.Regexp
	shortCut       string
	entries        StringList
}

func CITE(shortCut, filename, p string) *Cite {
	cite := Cite{
		commandPattern: regexp.MustCompile(fmt.Sprintf("^[^ ]+ PRIVMSG ([^ ]+) :%s%s.*$", p, shortCut)),
		indexPattern:   regexp.MustCompile(fmt.Sprintf("^[^ ]+ PRIVMSG ([^ ]+) :%s%s #([0-9]+)", p, shortCut)),
		processPattern: regexp.MustCompile(fmt.Sprintf("^[^ ]+ PRIVMSG ([^ ]+) :%s%s *([^ ]*)$", p, shortCut)),
		shortCut:       shortCut,
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
	return this.commandPattern.MatchString(line)
}

func (this *Cite) Process(line string, gossip *Gossip) {
	channel, response := this.process(line)
	gossip.SendMessage(channel, response)
}

func (this *Cite) process(line string) (channel, response string) {
	var (
		max, index int
		entry      string
		err        error
	)
	indexMatch := this.indexPattern.FindStringSubmatch(line)
	if indexMatch != nil {
		channel = indexMatch[1]
		max, index, entry, err = this.entries.Index(indexMatch[2])
		if err != nil {
			response = fmt.Sprintf("%s: %s", this.shortCut, err.Error())
		} else {
			response = fmt.Sprintf("%s[%d/%d]: %s", this.shortCut, index, max, entry)
		}
		return
	}
	match := this.processPattern.FindStringSubmatch(line)
	channel = match[1]
	filter := match[2]
	if filter != "" {
		max, index, entry = this.entries.FilteredRandom(filter)
	} else {
		max, index, entry = this.entries.Random()
	}
	response = fmt.Sprintf("%s[%d/%d]: %s", this.shortCut, index, max, entry)
	return
}
