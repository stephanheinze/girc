package main

import (
	"fmt"
	"log"
	"net/textproto"
)

type Command interface {
	Match(line string) bool
	Process(line string, gossip *Gossip)
}

type Gossip struct {
	Server   string
	Port     uint
	Channel  string
	Nick     string
	Conn     *textproto.Conn
	Commands []Command
}

func (this *Gossip) addCommand(command Command) {
	this.Commands = append(this.Commands, command)
}

func (this *Gossip) parseLine(line string) {
	log.Printf("%s", line)
	for _, command := range this.Commands {
		if command.Match(line) {
			command.Process(line, this)
		}
	}
}

func (this *Gossip) start() {
	c, err := textproto.Dial("tcp", fmt.Sprintf("%s:%d", this.Server, this.Port))
	if err != nil {
		log.Fatalf("Could not connect to server - reason: %s", err.Error())
		return
	}
	this.Conn = c
	c.Cmd("NICK %s\r\n", this.Nick)
	c.Cmd("USER %s 8 * :%s", this.Nick, this.Nick)
	if this.Channel != "" {
		c.Cmd("JOIN %s\r\n", this.Channel)
	}
	for {
		text, err := c.ReadLine()
		if err != nil {
			log.Fatalf("Could not read line - reason: %s", err.Error())
			return
		}
		go this.parseLine(text)
	}
}
