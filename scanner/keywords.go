package scanner

import "github.com/vishal-chdhry/lox.go/ast"

var keywords = map[string]ast.TokenType{
	"and":    ast.AND,
	"class":  ast.CLASS,
	"else":   ast.ELSE,
	"false":  ast.FALSE,
	"for":    ast.FOR,
	"fun":    ast.FUN,
	"if":     ast.IF,
	"nil":    ast.NIL,
	"or":     ast.OR,
	"print":  ast.PRINT,
	"return": ast.RETURN,
	"super":  ast.SUPER,
	"this":   ast.THIS,
	"true":   ast.TRUE,
	"var":    ast.VAR,
	"while":  ast.WHILE,
}
