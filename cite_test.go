package main

import "testing"

func TestEmptyCiteFile(t *testing.T) {
	cite := CITE("test", "examples/empty.txt")
	_, response := cite.process(":tester!tester@test.irc.server.org PRIVMSG #bla :#test")
	if response != "test[0/0]: -- no entries --" {
		t.Errorf("Invalid random cite on empty database. %q", response)
	}
}

func TestFilter(t *testing.T) {
	cite := CITE("test", "examples/cites.txt")
	_, response := cite.process(":tester!tester@test.irc.server.org PRIVMSG #bla :#test another")
	if response != "test[1/1]: another one" {
		t.Errorf("Invalid filtered cite. %q", response)
	}
}

func TestFilterNoResult(t *testing.T) {
	cite := CITE("test", "examples/cites.txt")
	_, response := cite.process(":tester!tester@test.irc.server.org PRIVMSG #bla :#test unknown")
	if response != "test[0/0]: -- no entries --" {
		t.Errorf("Invalid filtered response without result. %q", response)
	}
}

func TestCiteIndex(t *testing.T) {
	cite := CITE("test", "examples/cites.txt")
	_, response := cite.process(":tester!tester@test.irc.server.org PRIVMSG #bla :#test #8")
	if response != "test[8/8]: lorem ipsum" {
		t.Errorf("Invalid indexed cite. %q", response)
	}
}
