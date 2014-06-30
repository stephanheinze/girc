package main

import (
	"fmt"
	"regexp"
)

type Calc struct {
	pattern *regexp.Regexp
}

func CALC(p string) *Calc {
	return &Calc{
		pattern: regexp.MustCompile(fmt.Sprintf("^[^ ]+ PRIVMSG ([^ ]+) :%scalc(.*)$", p)),
	}
}

func (this *Calc) Match(line string) bool {
	return this.pattern.MatchString(line)
}

func (this *Calc) Process(line string, gossip *Gossip) {
	channel, response := this.process(line)
	gossip.SendMessage(channel, response)
}

func (this *Calc) process(line string) (string, string) {
	match := this.pattern.FindStringSubmatch(line)
	channel := match[1]
	equation := match[2]
	if len(equation) > 0 {
		result, err := this.calculate(equation)
		if err == nil {
			return channel, fmt.Sprintf("%s = %s\r\n", equation, result)
		}
		return channel, fmt.Sprintf("tsk ... %s\r\n", err.Error())
	}
	return channel, "no idea"
}

func (this *Calc) calculate(line string) (string, error) {
	return CalcParse(&CalcLex{s: line})
}
