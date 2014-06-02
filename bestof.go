package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"sync"
)

type BestOf struct {
	commandPattern *regexp.Regexp
	processPattern *regexp.Regexp
	filename       string
	store          bool
	entries        []string
	mutex          *sync.Mutex
}

func BEST_OF(filename string) *BestOf {
	bestOf := BestOf{
		commandPattern: regexp.MustCompile("^[^ ]+ PRIVMSG ([^ ]+) :!bestof.*$"),
		processPattern: regexp.MustCompile("^[^ ]+ PRIVMSG ([^ ]+) :!bestof([^ ]*)(.*)$"),
		filename:       filename,
		mutex:          &sync.Mutex{},
	}
	if filename == "" {
		bestOf.store = false
	} else {
		if file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666); err != nil {
			log.Printf("Could not open bestof-file %q - reason: %s", filename)
			log.Printf("!!! BestOfs will *NOT BE STORED* !!!")
			bestOf.store = false
		} else {
			defer file.Close()
			bestOf.store = true
			if err := json.NewDecoder(file).Decode(&bestOf.entries); err != nil {
				log.Printf("Could not decode bestof-file %q - reason: %s", filename, err.Error())
			}
		}
	}
	return &bestOf
}

func (this *BestOf) Match(line string) bool {
	return this.commandPattern.MatchString(line)
}

func (this *BestOf) Process(line string, gossip *Gossip) {
	channel, response := this.process(line)
	gossip.Conn.Cmd("PRIVMSG %s :%s\r\n", channel, response)
}

func (this *BestOf) process(line string) (channel string, response string) {
	match := this.processPattern.FindStringSubmatch(line)
	channel = match[1]
	switch match[2] {
	case "":
		max, index, entry := this.Random()
		response = fmt.Sprintf("BestOf[%d/%d]: %s", index, max, entry)
	case "-add":
		total := this.Add(match[3])
		if total == 1 {
			response = "Ok. Added. Got 1 bestof entry now."
		} else {
			response = fmt.Sprintf("Ok. Added. Got %d bestof entries now.", total)
		}
	case "-del":
		response = "not implemented yet."
	case "-count":
		response = fmt.Sprintf("I've got %d bestof entries.", len(this.entries))
	default:
		response = "unknown subcommand %q - use bestoff|bestof-add|bestof-count."
	}
	return
}

func (this *BestOf) Add(entry string) int {
	this.entries = append(this.entries, entry)
	if this.store {
		this.mutex.Lock()
		defer this.mutex.Unlock()
		if file, err := os.OpenFile(this.filename, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0666); err == nil {
			if err := json.NewEncoder(file).Encode(this.entries); err != nil {
				log.Printf("Can't encode bestof entries - reason: %s", err.Error())
			}
			file.Close()
		} else {
			log.Printf("Can't write bestof file - reason: %s", err.Error())
		}
	}
	return len(this.entries)
}

func (this *BestOf) Random() (max int, index int, entry string) {
	max = len(this.entries)
	switch max {
	case 0:
		return 0, 0, "-- no entries yet --"
	case 1:
		return 1, 1, this.entries[0]
	default:
		index = rand.Intn(max)
		entry = this.entries[index]
		index = index + 1
		return
	}
}
