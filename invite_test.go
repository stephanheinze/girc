package main

import (
	"testing"
)

func TestInviteMatchLineInvite(t *testing.T) {
	me := INVITE()
	b := me.Match(":tester!tester@test.irc.server.org INVITE girc :#channel2")
	if b == false {
		t.Error("Match wrong result.")
	}
	b = me.Match(":tester!tester@test.irc.server.org INVITE girc #channel2")
	if b == false {
		t.Error("Match wrong result.")
	}
}

func TestInviteMatchLinePrivMsg(t *testing.T) {
	me := INVITE()
	b := me.Match(":tester!tester@test.irc.server.org PRIVMSG #bla :what's up guys")
	if b == true {
		t.Error("Match wrong result.")
	}
}

func TestInviteParseLine(t *testing.T) {
	me := INVITE()
	channel, name := me.parseLine(":tester!tester@test.irc.server.org INVITE girc :#channel2")
	if channel != "#channel2" {
		t.Error("Wrong value for channel.")
	}
	if name != "tester" {
		t.Error("Wrong value for username.")
	}
}
