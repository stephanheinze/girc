package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/textproto"
	"os"
	"time"
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
	for {
		text, err := c.ReadLine()
		if err != nil {
			log.Fatalf("Could not read line - reason: %s", err.Error())
			return
		}
		go this.parseLine(text)
	}
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	portPtr := flag.Uint("p", uint(6667), "port")
	channelPtr := flag.String("c", "", "channel to join at startup")
	bestofPtr := flag.String("b", "/tmp/girc-bestof.json", "bestof-data-file")

	flag.Parse()

	gossip := new(Gossip)
	gossip.Channel = *channelPtr
	gossip.Port = *portPtr
	switch flag.NArg() {
	case 2:
		gossip.Nick = flag.Arg(0)
		gossip.Server = flag.Arg(1)
	case 1:
		gossip.Nick = "girc"
		gossip.Server = flag.Arg(0)
	default:
		fmt.Printf("server missing.\n")
		os.Exit(1)
	}

	gossip.addCommand(JOIN())
	gossip.addCommand(PING())

	gossip.addCommand(CALC())
	gossip.addCommand(KICK())
	gossip.addCommand(BEST_OF(*bestofPtr))

	gossip.start()
}
