package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	portPtr := flag.Uint("p", uint(6667), "port")
	channelPtr := flag.String("c", "", "channel to join at startup")
	bestofPtr := flag.String("b", "/tmp/girc-bestof.json", "bestof-data-file")
	mePtr := flag.String("me", "", "me-data-file")
	matrixPtr := flag.String("m", "", "file with one liner cites from matrix (the film)")
	cmdPrefixPtr := flag.String("t", "!", "prefix to trigger bot commands")

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

	gossip.addCommand(CALC(*cmdPrefixPtr))
	gossip.addCommand(KICK())
	gossip.addCommand(BEST_OF(*bestofPtr, *cmdPrefixPtr))
	gossip.addCommand(ME(gossip.Nick, *mePtr))
	gossip.addCommand(INVITE())
	gossip.addCommand(LEAVE(*cmdPrefixPtr))

	if *matrixPtr != "" {
		gossip.addCommand(CITE("matrix", *matrixPtr, *cmdPrefixPtr))
	}

	gossip.start()
}
