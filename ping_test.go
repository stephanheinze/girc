package main

import (
	"testing"
)

func TestMatch(t *testing.T) {
	me := PING()
	b := me.Match("PING :my-domain.my-host.com")
	if b == false {
		t.Errorf("Wrong match value.")
	}
}
