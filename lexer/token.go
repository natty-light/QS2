package lexer

type TokenType string

const (
	// Literals
	Null       TokenType = "Null"
	Identifier TokenType = "Identifier"
	Int        TokenType = "Int"

	// Keywords
	Mut    TokenType = "Mut"
	Const  TokenType = "Const"
	True   TokenType = "True"
	False  TokenType = "False"
	If     TokenType = "If"
	Else   TokenType = "Else"
	Elseif TokenType = "Elseif"
	Func   TokenType = "Func"
	Return TokenType = "Return"

	// Grouping
	LeftParen          TokenType = "LeftParen"
	RightParen         TokenType = "RightParen"
	LeftCurlyBracket   TokenType = "LeftCurlyBracket"
	RightCurlyBracket  TokenType = "RightCurlyBracket"
	LeftSquareBracket  TokenType = "LeftSquareBracket"
	RightSquareBracket TokenType = "RightSquareBracket"
	Semicolon          TokenType = "Semicolon"
	Comma              TokenType = "Comma"
	Colon              TokenType = "Colon"
	Dot                TokenType = "Dot"

	// Symbols
	Plus     TokenType = "Plus"
	Minus    TokenType = "Minus"
	Quotient TokenType = "Quotient"
	Product  TokenType = "Product"
	Modulo   TokenType = "Modulus"
	// Assign             TokenType = "Assign"
	Equality    TokenType = "Equality"
	GreaterThan TokenType = "GreaterThan"
	LessThan    TokenType = "LessThan"
	Band        TokenType = "Bang"
	// GreaterThanEqualTo TokenType = "GreaterThanEqualTo"
	// LessThanEqualTo    TokenType = "LessThanEqualTo"
	// NotEqual           TokenType = "NotEqual"
	And TokenType = "And"
	Or  TokenType = "Or"

	EOF     TokenType = "EOF" // End of File
	Illegal TokenType = "Illegal"
)

type Token struct {
	Literal string
	Type    TokenType
	Line    int
}

func token(Type TokenType, char byte, Line int) Token {
	return Token{Type: Type, Literal: string(char), Line: Line}
}
