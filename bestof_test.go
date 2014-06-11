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
	bestof := exampleBestOf()
	_, response := bestof.process(":tester!tester@test.irc.server.org PRIVMSG #bla :!bestof  ")
	if strings.HasPrefix("BestOf[", response) {
		t.Errorf("Invalid bestof response. %q", response)
	}
}

func TestBestOfFiltered(t *testing.T) {
	bestof := exampleBestOf()
	_, response := bestof.process("tester!tester@tester.irc.server.org PRIVMSG #bla :!bestof thr")
	if response != "BestOf[1/1]: entry three" {
		t.Errorf("Invalid filtered bestof response. %q", response)
	}
}

func TestBestOfIndexed(t *testing.T) {
	bestof := exampleBestOf()
	_, response := bestof.process("tester!tester@tester.irc.server.org PRIVMSG #bla :!bestof #2")
	if response != "Bestof[2/4]: entry two" {
		t.Errorf("Invalid indexed bestof response. %q", response)
	}
}

func TestBestOfIndexTooLarge(t *testing.T) {
	bestof := exampleBestOf()
	_, response := bestof.process("tester!tester@tester.irc.server.org PRIVMSG #bla :!bestof #5")
	if response != "Bestof: invalid index #5" {
		t.Errorf("Invalid bestof response for too large index. %q", response)
	}
}

func TestBestOfIndexZero(t *testing.T) {
	bestof := exampleBestOf()
	_, response := bestof.process("tester!tester@tester.irc.server.org PRIVMSG #bla :!bestof #0")
	if response != "Bestof: invalid index #0" {
		t.Errorf("Invalid bestof response for too large index. %q", response)
	}
}

func exampleBestOf() *BestOf {
	bestof := BEST_OF("")
	bestof.Add("entry one")
	bestof.Add("entry two")
	bestof.Add("entry three")
	bestof.Add("entry four")
	return bestof
}
