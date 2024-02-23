package lexer

import "testing"

func TestNextToken(t *testing.T) {
	source := `mut five = 5;
	mut ten = 10;
	mut add = func (x, y) {
		 x + y;
	};
	mut result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
		expectedLine    int
	}{
		{Mut, "mut", 1},
		{Identifier, "five", 1},
		{Assign, "=", 1},
		{Int, "5", 1},
		{Semicolon, ";", 1},

		{Mut, "mut", 2},
		{Identifier, "ten", 2},
		{Assign, "=", 2},
		{Int, "10", 2},
		{Semicolon, ";", 2},

		{Mut, "mut", 3},
		{Identifier, "add", 3},
		{Assign, "=", 3},
		{Func, "func", 3},
		{LeftParen, "(", 3},
		{Identifier, "x", 3},
		{Comma, ",", 3},
		{Identifier, "y", 3},
		{RightParen, ")", 3},
		{LeftCurlyBracket, "{", 3},

		{Identifier, "x", 4},
		{Plus, "+", 4},
		{Identifier, "y", 4},
		{Semicolon, ";", 4},

		{RightCurlyBracket, "}", 5},
		{Semicolon, ";", 5},

		{Mut, "mut", 6},
		{Identifier, "result", 6},
		{Assign, "=", 6},
		{Identifier, "add", 6},
		{LeftParen, "(", 6},
		{Identifier, "five", 6},
		{Comma, ",", 6},
		{Identifier, "ten", 6},
		{RightParen, ")", 6},
		{Semicolon, ";", 6},

		{Bang, "!", 7},
		{Minus, "-", 7},
		{Slash, "/", 7},
		{Star, "*", 7},
		{Int, "5", 7},
		{Semicolon, ";", 7},

		{Int, "5", 8},
		{LessThan, "<", 8},
		{Int, "10", 8},
		{GreaterThan, ">", 8},
		{Int, "5", 8},
		{Semicolon, ";", 8},

		{EOF, "", 0},
	}

	lexer := CreateLexer(source)

	for i, tt := range tests {
		tok := lexer.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d", i, tt.expectedLine, tok.Line)
		}
	}

}
