package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type Calc struct {
	commandPattern *regexp.Regexp
	processPattern *regexp.Regexp
}

func CALC() *Calc {
	return &Calc{
		commandPattern: regexp.MustCompile("^[^ ]+ PRIVMSG ([^ ]+) :!calc *.*$"),
		processPattern: regexp.MustCompile("^[^ ]+ PRIVMSG ([^ ]+) :!calc ([\\d]+) *([\\+\\-\\*\\:]) *([\\d]+).*$"),
	}
}

func (this *Calc) Match(line string) bool {
	return this.commandPattern.MatchString(line)
}

func (this *Calc) Process(line string, gossip *Gossip) {
	channel, response := this.process(line)
	gossip.Conn.Cmd("PRIVMSG %s :%s\r\n", channel, response)
}

func (this *Calc) process(line string) (string, string) {
	match := this.processPattern.FindStringSubmatch(line)
	if len(match) != 5 {
		return match[1], "expect me to answer that?"
	}
	var result string
	op1, _ := strconv.ParseUint(match[2], 10, 64)
	op2, _ := strconv.ParseUint(match[4], 10, 64)
	switch match[3] {
	case "+":
		result = fmt.Sprintf("%d", op1+op2)
	case "-":
		result = fmt.Sprintf("%d", op1-op2)
	case "*":
		result = fmt.Sprintf("%d", op1*op2)
	case ":":
		if op2 == 0 {
			result = "NaN"
		} else {
			remain := op1 % op2
			if remain != 0 {
				result = fmt.Sprintf("%d (rest %d)", op1/op2, remain)
			} else {
				result = fmt.Sprintf("%d", op1/op2)
			}
		}
	}
	return match[1], fmt.Sprintf(" %d %s %d = %s\r\n", op1, match[3], op2, result)
}
