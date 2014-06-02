package main

import (
	"strings"
	"testing"
)

func TestBestOfCountEmpty(t *testing.T) {
	bestof := BEST_OF("")
	_, response := bestof.process(":tester!tester@test.irc.server.org PRIVMSG #bla :!bestof-count")
	if response != "I've got 0 bestof entries." {
		t.Errorf("Invalid count on empty bestof. %q", response)
	}
}

func TestBestOfAdd(t *testing.T) {
	var response string
	bestof := BEST_OF("")
	_, response = bestof.process(":tester!tester@test.irc.server.org PRIVMSG #bla :!bestof-add hello world")
	if response != "Ok. Added. Got 1 bestof entry now." {
		t.Errorf("Invalid response on first bestof-add. %q", response)
	}
	_, response = bestof.process(":tester!tester@test.irc.server.org PRIVMSG #bla :!bestof-add hello outer world")
	if response != "Ok. Added. Got 2 bestof entries now." {
		t.Errorf("Invalid response on second bestof-add. %q", response)
	}
}

func TestBestOfRandom(t *testing.T) {
	bestof := BEST_OF("")
	bestof.Add("entry 1")
	bestof.Add("entry 2")
	bestof.Add("entry 3")
	bestof.Add("entry 4")
	_, response := bestof.process(":tester!tester@test.irc.server.org PRIVMSG #bla :!bestof  ")
	if strings.HasPrefix("BestOf[", response) {
		t.Errorf("Invalid bestof response. %q", response)
	}
}
