package main

import "testing"

func TestCalcCommandMatch(t *testing.T) {
	if !CALC().Match(":tester!tester@test.irc.server.org PRIVMSG #bla :!calc") {
		t.Error("pattern not matching as expected: !calc")
	}
	if !CALC().Match(":tester!tester@test.irc.server.org PRIVMSG #bla :!calc ") {
		t.Error("pattern not matching as expected. !calc ")
	}
	if !CALC().Match(":tester!tester@test.irc.server.org PRIVMSG #bla :!calc 2 + 2") {
		t.Error("pattern not matching as expected. !calc 2 + 2")
	}
	if !CALC().Match(":tester!tester@test.irc.server.org PRIVMSG #bla :!calc 2+2") {
		t.Error("pattern not matching as expected. !calc 2+2")
	}
}

func TestCalcProcess(t *testing.T) {
	calc := CALC()
	channel, result := calc.process(":tester!tester@test.irc.server.org PRIVMSG #bla :!calc 2 + 2")
	if channel != "#bla" {
		t.Errorf("invalid channel %q", channel)
	}
	if result != " 2 + 2 = 4\r\n" {
		t.Errorf("invalid result %q", result)
	}
}
