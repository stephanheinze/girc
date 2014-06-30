package main

import (
	"testing"
)

func TestLeaveMatchLineLeave(t *testing.T) {
	me := LEAVE("#")
	b := me.Match(":tester!tester@test.irc.server.org PRIVMSG #channel2 :#kick")
	if b == false {
		t.Error("Match wrong result.")
	}
}

func TestLeaveMatchLinePrivMsg(t *testing.T) {
	me := LEAVE("#")
	b := me.Match(":tester!tester@test.irc.server.org PRIVMSG #bla :what's up guys")
	if b == true {
		t.Error("Match wrong result.")
	}
}

func TestLeaveParseLine(t *testing.T) {
	me := LEAVE("#")
	channel, name := me.parseLine(":tester!tester@test.irc.server.org PRIVMSG #channel2 :#kick")
	if channel != "#channel2" {
		t.Error("Wrong value for channel.")
	}
	if name != "tester" {
		t.Error("Wrong value for username.")
	}
}
