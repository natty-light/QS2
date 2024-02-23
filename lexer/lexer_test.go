package lexer

import (
	"QuonkScript/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	source := `mut five = 5;
	mut ten = 10;
	mut add = func (x, y) {
		 x + y;
	};
	mut result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	if (5 < 10) {
		return true;
	} else {
		return false;
	}
	10 == 10;
	10 != 9;
	5 >= 10;
	7 <= 6;
	true && false;
	true || false;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
	}{
		{token.Mut, "mut", 1},
		{token.Identifier, "five", 1},
		{token.Assign, "=", 1},
		{token.Number, "5", 1},
		{token.Semicolon, ";", 1},

		{token.Mut, "mut", 2},
		{token.Identifier, "ten", 2},
		{token.Assign, "=", 2},
		{token.Number, "10", 2},
		{token.Semicolon, ";", 2},

		{token.Mut, "mut", 3},
		{token.Identifier, "add", 3},
		{token.Assign, "=", 3},
		{token.Func, "func", 3},
		{token.LeftParen, "(", 3},
		{token.Identifier, "x", 3},
		{token.Comma, ",", 3},
		{token.Identifier, "y", 3},
		{token.RightParen, ")", 3},
		{token.LeftCurlyBracket, "{", 3},

		{token.Identifier, "x", 4},
		{token.Plus, "+", 4},
		{token.Identifier, "y", 4},
		{token.Semicolon, ";", 4},

		{token.RightCurlyBracket, "}", 5},
		{token.Semicolon, ";", 5},

		{token.Mut, "mut", 6},
		{token.Identifier, "result", 6},
		{token.Assign, "=", 6},
		{token.Identifier, "add", 6},
		{token.LeftParen, "(", 6},
		{token.Identifier, "five", 6},
		{token.Comma, ",", 6},
		{token.Identifier, "ten", 6},
		{token.RightParen, ")", 6},
		{token.Semicolon, ";", 6},

		{token.Bang, "!", 7},
		{token.Minus, "-", 7},
		{token.Slash, "/", 7},
		{token.Star, "*", 7},
		{token.Number, "5", 7},
		{token.Semicolon, ";", 7},

		{token.Number, "5", 8},
		{token.LessThan, "<", 8},
		{token.Number, "10", 8},
		{token.GreaterThan, ">", 8},
		{token.Number, "5", 8},
		{token.Semicolon, ";", 8},

		{token.If, "if", 9},
		{token.LeftParen, "(", 9},
		{token.Number, "5", 9},
		{token.LessThan, "<", 9},
		{token.Number, "10", 9},
		{token.RightParen, ")", 9},
		{token.LeftCurlyBracket, "{", 9},

		{token.Return, "return", 10},
		{token.True, "true", 10},
		{token.Semicolon, ";", 10},

		{token.RightCurlyBracket, "}", 11},
		{token.Else, "else", 11},
		{token.LeftCurlyBracket, "{", 11},

		{token.Return, "return", 12},
		{token.False, "false", 12},
		{token.Semicolon, ";", 12},

		{token.RightCurlyBracket, "}", 13},

		{token.Number, "10", 14},
		{token.EqualTo, "==", 14},
		{token.Number, "10", 14},
		{token.Semicolon, ";", 14},

		{token.Number, "10", 15},
		{token.NotEqualTo, "!=", 15},
		{token.Number, "9", 15},
		{token.Semicolon, ";", 15},

		{token.Number, "5", 16},
		{token.GreaterThanEqualTo, ">=", 16},
		{token.Number, "10", 16},
		{token.Semicolon, ";", 16},

		{token.Number, "7", 17},
		{token.LessThanEqualTo, "<=", 17},
		{token.Number, "6", 17},
		{token.Semicolon, ";", 17},

		{token.True, "true", 18},
		{token.And, "&&", 18},
		{token.False, "false", 18},
		{token.Semicolon, ";", 18},

		{token.True, "true", 19},
		{token.Or, "||", 19},
		{token.False, "false", 19},
		{token.Semicolon, ";", 19},

		{token.EOF, "", 0},
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
