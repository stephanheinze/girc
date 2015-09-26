package main

import (
	"strings"
	"testing"
)

func TestExtListCountEmpty(t *testing.T) {
	mylist := EXT_LIST("mylist", "", "#")
	_, response := mylist.process(":tester!tester@test.irc.server.org PRIVMSG #bla :#mylist-count")
	if response != "I've got 0 mylist entries." {
		t.Errorf("Invalid count on empty mylist. %q", response)
	}
}

func TestExtListAdd(t *testing.T) {
	mylist := EXT_LIST("mylist", "", "#")
	_, response := mylist.process(":tester!tester@test.irc.server.org PRIVMSG #bla :#mylist-add hello world")
	if response != "OK" {
		t.Errorf("Invalid response on first mylist-add. %q", response)
	}
	_, response = mylist.process(":tester!tester@test.irc.server.org PRIVMSG #bla :#mylist-add hello outer world")
	if response != "OK" {
		t.Errorf("Invalid response on second mylist-add. %q", response)
	}
}

func TestExtListAddNothing(t *testing.T) {
	mylist := EXT_LIST("mylist", "", "#")
	_, response := mylist.process(":tester!tester@test.irc.server.org PRIVMSG #bla :#mylist-add")
	if response != "error: nothing to add" {
		t.Errorf("Invalid response for empty mylist-add. %q", response)
	}
}

func TestExtListAddSpace(t *testing.T) {
	mylist := EXT_LIST("mylist", "", "#")
	_, response := mylist.process(":tester!tester@test.irc.server.org PRIVMSG #bla :#mylist-add   ")
	if response != "error: nothing to add" {
		t.Errorf("Invalid response for empty mylist-add. %q", response)
	}
}

func TestExtListRandom(t *testing.T) {
	mylist := exampleBestOf()
	_, response := mylist.process(":tester!tester@test.irc.server.org PRIVMSG #bla :#mylist  ")
	if strings.HasPrefix("mylist[", response) {
		t.Errorf("Invalid mylist response. %q", response)
	}
}

func TestExtListFiltered(t *testing.T) {
	mylist := exampleBestOf()
	_, response := mylist.process("tester!tester@tester.irc.server.org PRIVMSG #bla :#mylist thr")
	if response != "mylist[1/1]: entry three" {
		t.Errorf("Invalid filtered mylist response. %q", response)
	}
}

func TestExtListIndexed(t *testing.T) {
	mylist := exampleBestOf()
	_, response := mylist.process("tester!tester@tester.irc.server.org PRIVMSG #bla :#mylist #2")
	if response != "mylist[2/4]: entry two" {
		t.Errorf("Invalid indexed mylist response. %q", response)
	}
}

func TestExtListIndexTooLarge(t *testing.T) {
	mylist := exampleBestOf()
	_, response := mylist.process("tester!tester@tester.irc.server.org PRIVMSG #bla :#mylist #5")
	if response != "mylist: invalid index #5" {
		t.Errorf("Invalid mylist response for too large index. %q", response)
	}
}

func TestExtListIndexZero(t *testing.T) {
	mylist := exampleBestOf()
	_, response := mylist.process("tester!tester@tester.irc.server.org PRIVMSG #bla :#mylist #0")
	if response != "mylist: invalid index #0" {
		t.Errorf("Invalid mylist response for too large index. %q", response)
	}
}

func exampleBestOf() *ExtList {
	mylist := EXT_LIST("mylist", "", "#")
	mylist.Add("entry one")
	mylist.Add("entry two")
	mylist.Add("entry three")
	mylist.Add("entry four")
	return mylist
}
