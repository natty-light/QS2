package lexer

type TokenType string

const (
	// Literals
	Null       TokenType = "Null"
	Identifier TokenType = "Identifier"
	Number     TokenType = "Number"

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
	Plus        TokenType = "Plus"
	Minus       TokenType = "Minus"
	Slash       TokenType = "Slash"
	Star        TokenType = "Star"
	Modulo      TokenType = "Modulus"
	Assign      TokenType = "Assign"
	GreaterThan TokenType = "GreaterThan"
	LessThan    TokenType = "LessThan"
	Bang        TokenType = "Bang"

	// Multi char symbols
	EqualTo            TokenType = "Equality"
	GreaterThanEqualTo TokenType = "GreaterThanEqualTo"
	LessThanEqualTo    TokenType = "LessThanEqualTo"
	NotEqualTo         TokenType = "NotEqual"
	And                TokenType = "And"
	Or                 TokenType = "Or"

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
