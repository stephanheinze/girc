package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

type ExtList struct {
	commandPattern *regexp.Regexp
	indexPattern   *regexp.Regexp
	processPattern *regexp.Regexp
	shortcut       string
	filename       string
	store          bool
	entries        StringList
	mutex          *sync.Mutex
}

func EXT_LIST(shortcut, filename, p string) *ExtList {
	bestOf := ExtList{
		commandPattern: regexp.MustCompile(fmt.Sprintf("^[^ ]+ PRIVMSG ([^ ]+) :%s%s.*$", p, shortcut)),
		indexPattern:   regexp.MustCompile(fmt.Sprintf("^[^ ]+ PRIVMSG ([^ ]+) :%s%s #([0-9]+)", p, shortcut)),
		processPattern: regexp.MustCompile(fmt.Sprintf("^[^ ]+ PRIVMSG ([^ ]+) :%s%s([^ ]*)(.*)$", p, shortcut)),
		shortcut:       shortcut,
		filename:       filename,
		mutex:          &sync.Mutex{},
	}
	if filename == "" {
		bestOf.store = false
	} else {
		if file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666); err != nil {
			log.Printf("Could not open %s-file %q - reason: %s", bestOf.shortcut, filename, err.Error())
			log.Printf("!!! New entries *WILL NOT BE STORED* !!!")
			bestOf.store = false
		} else {
			defer file.Close()
			bestOf.store = true
			bestOf.entries.Read(file)
		}
	}
	return &bestOf
}

func (this *ExtList) Match(line string) bool {
	return this.commandPattern.MatchString(line)
}

func (this *ExtList) Process(line string, gossip *Gossip) {
	channel, response := this.process(line)
	gossip.SendMessage(channel, response)
}

func (this *ExtList) process(line string) (channel string, response string) {
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
			response = fmt.Sprintf("%s: %s", this.shortcut, err.Error())
		} else {
			response = fmt.Sprintf("%s[%d/%d]: %s", this.shortcut, index, max, entry)
		}
		return
	}
	match := this.processPattern.FindStringSubmatch(line)
	channel = match[1]
	switch match[2] {
	case "":
		if match[3] != "" {
			max, index, entry = this.entries.FilteredRandom(match[3])
		} else {
			max, index, entry = this.entries.Random()
		}
		response = fmt.Sprintf("%s[%d/%d]: %s", this.shortcut, index, max, entry)
	case "-add":
		if len(strings.TrimSpace(match[3])) != 0 {
			this.Add(match[3])
			response = "OK"
		} else {
			response = "error: nothing to add"
		}
	case "-del":
		response = "not implemented yet."
	case "-count":
		response = fmt.Sprintf("I've got %d %s entries.", this.entries.Len(), this.shortcut)
	default:
		response = fmt.Sprintf("unknown subcommand - use %s|%s-add|%s-count.", this.shortcut, this.shortcut, this.shortcut)
	}
	return
}

func (this *ExtList) Add(entry string) int {
	this.entries.Add(entry)
	if this.store {
		this.mutex.Lock()
		defer this.mutex.Unlock()
		if file, err := os.OpenFile(this.filename, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0666); err == nil {
			if err := this.entries.Write(file); err != nil {
				log.Printf("Can't write %s entries - reason: %s", this.shortcut, err.Error())
			}
			file.Close()
		} else {
			log.Printf("Can't write %s file - reason: %s", this.shortcut, err.Error())
		}
	}
	return this.entries.Len()
}
