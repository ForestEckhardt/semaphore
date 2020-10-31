package semaphore

import (
	"errors"
	"unicode/utf8"
)

type FlagParser struct {
	input string
	pos   int
	start int
	width int

	flags []string
}

func NewFlagParser() FlagParser {
	return FlagParser{}
}

const eof = -1

func (p *FlagParser) next() rune {
	if p.pos == len(p.input) {
		return eof
	}

	r, w := utf8.DecodeRuneInString(p.input[p.pos:])
	p.pos += w
	p.width = w

	return r
}

func (p *FlagParser) backup() {
	p.pos -= p.width
}

func (p *FlagParser) skip() {
	p.start = p.pos
}

func (p *FlagParser) append() {
	p.flags = append(p.flags, p.input[p.start:p.pos])
	p.start = p.pos
}

func (p *FlagParser) ParseFlags(input string) ([]string, error) {
	p.input = input

	var r rune
	for r != eof {
		r = p.next()
		switch r {
		case ' ':
			p.backup()
			p.append()
			r = p.next()
			for r == ' ' {
				p.skip()
				r = p.next()
				if r == eof {
					break
				}
			}

		case '"':
			r = p.next()
			for r != '"' {
				r = p.next()
				if r == eof {
					return nil, errors.New("expected closing \" before end of input")
				}
			}

		case '\'':
			r = p.next()
			for r != '\'' {
				r = p.next()
				if r == eof {
					return nil, errors.New("expected closing ' before end of input")
				}
			}

		case eof:
			p.append()
		}

	}

	return p.flags, nil
}
