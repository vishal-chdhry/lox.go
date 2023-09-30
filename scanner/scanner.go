package scanner

import (
	"strconv"

	"github.com/vishal-chdhry/lox.go/ast"
	"github.com/vishal-chdhry/lox.go/errors"
	"go.uber.org/multierr"
)

type Scanner interface {
	ScanTokens() ([]ast.Token, error)
}

type scanner struct {
	source  string
	start   int
	current int
	line    int
	tokens  []ast.Token
}

func NewScanner(src string) Scanner {
	return &scanner{
		source:  src,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *scanner) ScanTokens() ([]ast.Token, error) {
	for {
		if s.isAtEnd() {
			break
		}
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}
	s.tokens = append(s.tokens, ast.Token{TokenType: ast.EOF, Lexeme: "", Literal: "", Line: s.line})
	return s.tokens, nil
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if char := rune(s.source[s.current]); char != expected {
		return false
	}

	s.current++
	return true
}

func (s *scanner) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (s *scanner) isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '_'
}

func (s *scanner) isAlphaNumberic(r rune) bool {
	return s.isAlpha(r) || s.isDigit(r)
}

func (s *scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	} else {
		return rune(s.source[s.current])
	}
}

func (s *scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return rune(s.source[s.current+1])
}

func (s *scanner) scanToken() error {
	var errs []error
	r := s.advance()
	switch r {
	case '(':
		s.addToken(ast.LEFT_PAREN)
	case ')':
		s.addToken(ast.RIGHT_PAREN)
	case '{':
		s.addToken(ast.LEFT_BRACE)
	case '}':
		s.addToken(ast.RIGHT_BRACE)
	case ',':
		s.addToken(ast.COMMA)
	case '.':
		s.addToken(ast.DOT)
	case '-':
		s.addToken(ast.MINUS)
	case '+':
		s.addToken(ast.PLUS)
	case ';':
		s.addToken(ast.SEMICOLON)
	case '*':
		s.addToken(ast.STAR)
	case '!':
		if s.match('=') {
			s.addToken(ast.BANG_EQUAL)
		} else {
			s.addToken(ast.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(ast.EQUAL_EQUAL)
		} else {
			s.addToken(ast.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(ast.LESS_EQUAL)
		} else {
			s.addToken(ast.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(ast.GREATER_EQUAL)
		} else {
			s.addToken(ast.GREATER)
		}
	case '/':
		if s.match('/') {
			for {
				if s.peek() == '\n' || s.isAtEnd() {
					break
				}
				s.advance()
			}
		} else {
			s.addToken(ast.SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
	case '\n':
		s.line++
	case '"':
		err := s.string()
		if err != nil {
			errs = append(errs, err)
		}
	default:
		var err error
		if s.isDigit(r) {
			err = s.number()
		} else if s.isAlpha(r) {
			err = s.identifier()
		} else {
			err = errors.Error(s.line, "", "Unexpected character.")
		}

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		err := multierr.Combine(errs...)
		return err
	}
	return nil
}

func (s *scanner) advance() rune {
	char := rune(s.source[s.current])
	s.current++
	return char
}

func (s *scanner) addToken(token ast.TokenType) {
	s.addTokenWithLiteral(token, nil)
}

func (s *scanner) addTokenWithLiteral(tokenType ast.TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	token := ast.Token{
		TokenType: tokenType,
		Lexeme:    text,
		Literal:   literal,
		Line:      s.line,
	}
	s.tokens = append(s.tokens, token)
}

func (s *scanner) string() error {
	for {
		if s.peek() == '"' || s.isAtEnd() {
			break
		}

		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if !s.isAtEnd() {
		return errors.Error(s.line, "", "Unterminated string")
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(ast.STRING, value)
	return nil
}

func (s *scanner) number() error {
	for {
		if !s.isDigit(s.peek()) {
			break
		}
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for {
			if !s.isDigit(s.peek()) {
				break
			}
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		return err
	}

	s.addTokenWithLiteral(ast.NUMBER, value)
	return nil
}

func (s *scanner) identifier() error {
	for {
		if !s.isAlphaNumberic(s.peek()) {
			break
		}
		s.advance()
	}
	text := s.source[s.start:s.current]
	if tokenType, ok := keywords[text]; ok {
		s.addToken(tokenType)
	} else {
		s.addToken(ast.IDENTIFIER)
	}
	return nil
}
