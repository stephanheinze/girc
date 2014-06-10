package main

import (
	"testing"
)

func TestMatchNotMentioned(t *testing.T) {
	me := ME("gossip", "")
	b := me.Match(":tester!tester@test.irc.server.org PRIVMSG #bla :what's up guys")
	if b == true {
		t.Errorf("Wrong match value.")
	}
}

func TestMatchNameBegin(t *testing.T) {
	me := ME("gossip", "")
	b := me.Match(":tester!tester@test.irc.server.org PRIVMSG #bla :gossip: what's up?")
	if b == false {
		t.Errorf("Wrong match value.")
	}
}

func TestMatchNameEnd(t *testing.T) {
	me := ME("gossip", "")
	b := me.Match(":tester!tester@test.irc.server.org PRIVMSG #bla :what's up, gossip?")
	if b == false {
		t.Errorf("Wrong match value.")
	}
}

func TestMatchNameMiddle(t *testing.T) {
	me := ME("gossip", "")
	b := me.Match(":tester!tester@test.irc.server.org PRIVMSG #bla :Hey gossip, what's up?")
	if b == false {
		t.Errorf("Wrong match value.")
	}
}

func TestProcessDefaultEntry(t *testing.T) {
	me := ME("gossip", "")
	_, response := me.process(":tester!tester@test.irc.server.org PRIVMSG #bla :hey gossip")
	if response != "Hey tester, stay calm! I'm just a poor bot." {
		t.Errorf("Invalid response: %q", response)
	}
}
