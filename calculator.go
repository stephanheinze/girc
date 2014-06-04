package main

import (
	"errors"
	"fmt"
	"unicode"
)

var regs = make([]int, 26)
var base int

type CalcSymType struct {
	yys int
	val int
}

const DIGIT = 57346
const LETTER = 57347
const UMINUS = 57348

var CalcToknames = []string{
	"DIGIT",
	"LETTER",
	" |",
	" &",
	" +",
	" -",
	" *",
	" /",
	" %",
	"UMINUS",
}
var CalcStatenames = []string{}

const CalcEofCode = 1
const CalcErrCode = 2
const CalcMaxDepth = 200

type CalcLex struct {
	s   string
	pos int
}

func (l *CalcLex) Lex(lval *CalcSymType) int {
	var c rune = ' '
	for c == ' ' {
		if l.pos == len(l.s) {
			return 0
		}
		c = rune(l.s[l.pos])
		l.pos += 1
	}

	if unicode.IsDigit(c) {
		lval.val = int(c - '0')
		return DIGIT
	} else if unicode.IsLower(c) {
		lval.val = int(c - 'a')
		return LETTER
	}
	return int(c)
}

func (l *CalcLex) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}

var CalcExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const CalcNprod = 18
const CalcPrivate = 57344

var CalcTokenNames []string
var CalcStates []string

const CalcLast = 54

var CalcAct = []int{

	3, 10, 11, 12, 13, 14, 18, 20, 21, 17,
	9, 22, 23, 24, 25, 26, 27, 28, 29, 16,
	15, 10, 11, 12, 13, 14, 8, 19, 8, 4,
	30, 6, 2, 6, 1, 12, 13, 14, 5, 7,
	5, 16, 15, 10, 11, 12, 13, 14, 15, 10,
	11, 12, 13, 14,
}
var CalcPact = []int{

	-1000, 24, -4, 35, -6, 22, 22, 4, -1000, -1000,
	22, 22, 22, 22, 22, 22, 22, 22, 13, -1000,
	-1000, -1000, 25, 25, -1000, -1000, -1000, -7, 41, 35,
	-1000,
}
var CalcPgo = []int{

	0, 0, 39, 34, 32,
}
var CalcR1 = []int{

	0, 3, 3, 4, 4, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 2, 2,
}
var CalcR2 = []int{

	0, 0, 3, 1, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 2, 1, 1, 1, 2,
}
var CalcChk = []int{

	-1000, -3, -4, -1, 5, 16, 9, -2, 4, 14,
	8, 9, 10, 11, 12, 7, 6, 15, -1, 5,
	-1, 4, -1, -1, -1, -1, -1, -1, -1, -1,
	17,
}
var CalcDef = []int{

	1, -2, 0, 3, 14, 0, 0, 15, 16, 2,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 14,
	13, 17, 6, 7, 8, 9, 10, 11, 12, 4,
	5,
}
var CalcTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	14, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 12, 7, 3,
	16, 17, 10, 8, 3, 9, 3, 11, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 15, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 6,
}
var CalcTok2 = []int{

	2, 3, 4, 5, 13,
}
var CalcTok3 = []int{
	0,
}

type CalcLexer interface {
	Lex(lval *CalcSymType) int
	Error(s string)
}

const CalcFlag = -1000

func CalcTokname(c int) string {
	if c >= 4 && c-4 < len(CalcToknames) {
		if CalcToknames[c-4] != "" {
			return CalcToknames[c-4]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func CalcStatname(s int) string {
	if s >= 0 && s < len(CalcStatenames) {
		if CalcStatenames[s] != "" {
			return CalcStatenames[s]
		}
	}
	return fmt.Sprintf("state-%v", s)
}

func Calclex1(lex CalcLexer, lval *CalcSymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = CalcTok1[0]
		goto out
	}
	if char < len(CalcTok1) {
		c = CalcTok1[char]
		goto out
	}
	if char >= CalcPrivate {
		if char < CalcPrivate+len(CalcTok2) {
			c = CalcTok2[char-CalcPrivate]
			goto out
		}
	}
	for i := 0; i < len(CalcTok3); i += 2 {
		c = CalcTok3[i+0]
		if c == char {
			c = CalcTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = CalcTok2[1]
	}
	return c
}

func CalcParse(Calclex CalcLexer) (string, error) {
	var Calcn int
	var Calclval CalcSymType
	var CalcVAL CalcSymType
	CalcS := make([]CalcSymType, CalcMaxDepth)

	Errflag := 0
	Calcstate := 0
	Calcchar := -1
	Calcp := -1
	goto Calcstack

ret0:
	return "hmpf 0", nil

ret1:
	return "hmpf 1", nil

Calcstack:

	Calcp++
	if Calcp >= len(CalcS) {
		nyys := make([]CalcSymType, len(CalcS)*2)
		copy(nyys, CalcS)
		CalcS = nyys
	}
	CalcS[Calcp] = CalcVAL
	CalcS[Calcp].yys = Calcstate

Calcnewstate:
	Calcn = CalcPact[Calcstate]
	if Calcn <= CalcFlag {
		goto Calcdefault
	}
	if Calcchar < 0 {
		Calcchar = Calclex1(Calclex, &Calclval)
	}
	Calcn += Calcchar
	if Calcn < 0 || Calcn >= CalcLast {
		goto Calcdefault
	}
	Calcn = CalcAct[Calcn]
	if CalcChk[Calcn] == Calcchar {
		Calcchar = -1
		CalcVAL = Calclval
		Calcstate = Calcn
		if Errflag > 0 {
			Errflag--
		}
		goto Calcstack
	}

Calcdefault:
	Calcn = CalcDef[Calcstate]
	if Calcn == -2 {
		if Calcchar < 0 {
			Calcchar = Calclex1(Calclex, &Calclval)
		}

		xi := 0
		for {
			if CalcExca[xi+0] == -1 && CalcExca[xi+1] == Calcstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			Calcn = CalcExca[xi+0]
			if Calcn < 0 || Calcn == Calcchar {
				break
			}
		}
		Calcn = CalcExca[xi+1]
		if Calcn < 0 {
			goto ret0
		}
	}
	if Calcn == 0 {
		switch Errflag {
		case 0:
			Calclex.Error("syntax error")
			fallthrough

		case 1, 2:
			Errflag = 3

			for Calcp >= 0 {
				Calcn = CalcPact[CalcS[Calcp].yys] + CalcErrCode
				if Calcn >= 0 && Calcn < CalcLast {
					Calcstate = CalcAct[Calcn]
					if CalcChk[Calcstate] == CalcErrCode {
						goto Calcstack
					}
				}
				Calcp--
			}
			goto ret1

		case 3:
			if Calcchar == CalcEofCode {
				goto ret1
			}
			Calcchar = -1
			goto Calcnewstate
		}
	}

	Calcnt := Calcn
	Calcpt := Calcp
	_ = Calcpt

	Calcp -= CalcR2[Calcn]
	CalcVAL = CalcS[Calcp+1]

	Calcn = CalcR1[Calcn]
	Calcg := CalcPgo[Calcn]
	Calcj := Calcg + CalcS[Calcp].yys + 1

	if Calcj >= CalcLast {
		Calcstate = CalcAct[Calcg]
	} else {
		Calcstate = CalcAct[Calcj]
		if CalcChk[Calcstate] != -Calcn {
			Calcstate = CalcAct[Calcg]
		}
	}
	switch Calcnt {

	case 3:
		return fmt.Sprintf("%d", CalcS[Calcpt-0].val), nil
	case 4:
		regs[CalcS[Calcpt-2].val] = CalcS[Calcpt-0].val
	case 5:
		CalcVAL.val = CalcS[Calcpt-1].val
	case 6:
		CalcVAL.val = CalcS[Calcpt-2].val + CalcS[Calcpt-0].val
	case 7:
		CalcVAL.val = CalcS[Calcpt-2].val - CalcS[Calcpt-0].val
	case 8:
		CalcVAL.val = CalcS[Calcpt-2].val * CalcS[Calcpt-0].val
	case 9:
		if CalcS[Calcpt-0].val == 0 {
			return "", errors.New("division by zero")
		}
		CalcVAL.val = CalcS[Calcpt-2].val / CalcS[Calcpt-0].val
	case 10:
		if CalcS[Calcpt-0].val == 0 {
			return "", errors.New("division by zero")
		}
		CalcVAL.val = CalcS[Calcpt-2].val % CalcS[Calcpt-0].val
	case 11:
		CalcVAL.val = CalcS[Calcpt-2].val & CalcS[Calcpt-0].val
	case 12:
		CalcVAL.val = CalcS[Calcpt-2].val | CalcS[Calcpt-0].val
	case 13:
		CalcVAL.val = -CalcS[Calcpt-0].val
	case 14:
		CalcVAL.val = regs[CalcS[Calcpt-0].val]
	case 15:
		CalcVAL.val = CalcS[Calcpt-0].val
	case 16:
		CalcVAL.val = CalcS[Calcpt-0].val
		if CalcS[Calcpt-0].val == 0 {
			base = 8
		} else {
			base = 10
		}
	case 17:
		CalcVAL.val = base*CalcS[Calcpt-1].val + CalcS[Calcpt-0].val
	}
	goto Calcstack
}
