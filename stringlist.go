package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

type StringList struct {
	entries []string
}

func (this *StringList) Add(entry string) {
	this.entries = append(this.entries, entry)
}

func (this *StringList) Read(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		this.Add(scanner.Text())
	}
	return scanner.Err()
}

func (this *StringList) Write(w io.Writer) error {
	for _, entry := range this.entries {
		if _, err := w.Write([]byte(fmt.Sprintf("%s\r\n", entry))); err != nil {
			return err
		}
	}
	return nil
}

func (this *StringList) Len() int {
	return len(this.entries)
}

func (this *StringList) Random() (max, index int, entry string) {
	return this.random(this.entries)
}

func (this *StringList) FilteredRandom(filter string) (max, index int, entry string) {
	reduced := make([]string, 0)
	loweredFilter := strings.ToLower(strings.TrimSpace(filter))
	for _, entry := range this.entries {
		if strings.Contains(strings.ToLower(entry), loweredFilter) {
			reduced = append(reduced, entry)
		}
	}
	return this.random(reduced)
}

func (this *StringList) random(entries []string) (max, index int, entry string) {
	max = len(entries)
	switch max {
	case 0:
		return 0, 0, "-- no entries --"
	case 1:
		return 1, 1, entries[0]
	default:
		index = rand.Intn(max)
		entry = entries[index]
		index = index + 1
		return
	}
}
