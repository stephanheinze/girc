package main

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"net/textproto"
	"time"
)

type Command interface {
	Match(line string) bool
	Process(line string, gossip *Gossip)
}

type Gossip struct {
	Server   string
	Port     uint
	Channels list.List
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

func (this *Gossip) SendMessage(channel, msg string) {
	this.Conn.Cmd("PRIVMSG %s :%s\r\n", channel, msg)
}

func (this *Gossip) JoinChannel(channel string) {
	this.Conn.Cmd("JOIN %s\r\n", channel)
}

func (this *Gossip) LeaveChannel(channel, reason string) {
	this.Conn.Cmd("PART %s :%s\r\n", channel, reason)
	this.RemoveChannel(channel)
}

func (this *Gossip) AddChannel(channel string) {
	this.Channels.PushBack(channel)
}

func (this *Gossip) RemoveChannel(channel string) {
	for e := this.Channels.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == channel {
			this.Channels.Remove(e)
			break
		}
	}
}

func (this *Gossip) PrintChannels() {
	fmt.Printf("%s is member of the following channels: ", this.Nick)
	for e := this.Channels.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value.(string))
		fmt.Print(" ")
	}
	fmt.Println()
}

func (this *Gossip) JoinChannels() {
	for e := this.Channels.Front(); e != nil; e = e.Next() {
		this.JoinChannel(e.Value.(string))
	}
}

func (this *Gossip) Connect() error {
	if this.Conn != nil {
		return errors.New("Already connected.")
	}

	c, err := textproto.Dial("tcp", fmt.Sprintf("%s:%d", this.Server, this.Port))
	if err != nil {
		return err
	}
	this.Conn = c

	c.Cmd("NICK %s\r\n", this.Nick)
	c.Cmd("USER %s 8 * :%s", this.Nick, this.Nick)

	return nil
}

func (this *Gossip) start() {
	if err := this.Connect(); err != nil {
		log.Fatalf("Could not connect to server - reason: %s", err.Error())
		return
	}

	for {
		text, err := this.Conn.ReadLine()
		if err != nil {
			log.Printf("Could not read line - reason: %s", err.Error())
			this.Conn = nil
			for {
				log.Println("Try to reconnect in 30 seconds...")
				time.Sleep(30 * time.Second)

				if err := this.Connect(); err != nil {
					log.Printf("Could not connect to server - reason: %s", err.Error())
					continue
				}

				this.JoinChannels()
				break
			}
			continue
		}

		go this.parseLine(text)

	}
}
