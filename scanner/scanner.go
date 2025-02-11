package scanner

import (
	"strconv"

	"github.com/drewslam/goloxTreeInterpreter/errors"
	"github.com/drewslam/goloxTreeInterpreter/token"
)

var keywords = map[string]token.TokenType{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

type Scanner struct {
	Source  string
	Tokens  []token.Token
	Start   int
	Current int
	Line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Source:  source,
		Tokens:  []token.Token{},
		Start:   0,
		Current: 0,
		Line:    1,
	}
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.Start = s.Current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, token.Token{
		Type:   token.EOF,
		Lexeme: "",
		Line:   s.Line,
	})
	return s.Tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
	case '{':
		s.addToken(token.LEFT_BRACE, nil)
	case '}':
		s.addToken(token.RIGHT_BRACE, nil)
	case ',':
		s.addToken(token.COMMA, nil)
	case '.':
		s.addToken(token.DOT, nil)
	case '-':
		s.addToken(token.MINUS, nil)
	case '+':
		s.addToken(token.PLUS, nil)
	case ';':
		s.addToken(token.SEMICOLON, nil)
	case '*':
		s.addToken(token.STAR, nil)
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL, nil)
		} else {
			s.addToken(token.BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL, nil)
		} else {
			s.addToken(token.EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL, nil)
		} else {
			s.addToken(token.LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL, nil)
		} else {
			s.addToken(token.GREATER, nil)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		s.Line++
	case '"':
		s.string()
	case 'o':
		if s.match('r') {
			s.addToken(token.OR, nil)
		}
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			errors.ReportError(s.Line, "Unexpected character.")
		}
	}
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.Source[s.Start:s.Current]
	tokenType, found := keywords[text]
	if !found {
		tokenType = token.IDENTIFIER
	}

	s.addToken(tokenType, nil)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Look for fractional part.
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the '.'
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	numStr := s.Source[s.Start:s.Current]
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		errors.ReportError(s.Line, "Invalid number.")
		return
	}
	s.addToken(token.NUMBER, num)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() != '\n' {
			s.Line++
		}
		s.advance()

		if s.isAtEnd() {
			errors.ReportError(s.Line, "Unterminated string.")
			return
		}

		// The closing "
		s.advance()

		// Trim the surrounding quotes.
		value := s.Source[s.Start+1 : s.Current-1]
		s.addToken(token.STRING, value)
	}
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.Source[s.Current] != expected {
		return false
	}

	s.Current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.Source[s.Current]
}

func (s *Scanner) peekNext() byte {
	if s.Current+1 >= len(s.Source) {
		return '\x00'
	}
	return s.Source[s.Current+1]
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) advance() byte {
	s.Current++
	return s.Source[s.Current-1]
}

func (s *Scanner) addToken(tokenType token.TokenType, literal interface{}) {
	lexeme := s.Source[s.Start:s.Current]

	s.Tokens = append(s.Tokens, token.Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    s.Line,
	})
}
