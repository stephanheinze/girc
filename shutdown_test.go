package main

import (
	"testing"
)

func TestNoMatchLineShutdown(t *testing.T) {
	me := SHUTDOWN("11111")
	b := me.Match(":tester!tester@test.irc.server.org PRIVMSG girc :shutdown 1111x")
	if b == true {
		t.Error("Match wrong result.")
	}
}

func TestMatchLineShotdown(t *testing.T) {
	me := SHUTDOWN("11111")
	b := me.Match(":tester!tester@test.irc.server.org PRIVMSG girc :shutdown 11111")
	if b == false {
		t.Error("Match wrong result.")
	}
}
