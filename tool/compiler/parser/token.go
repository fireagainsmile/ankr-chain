package parser

type Token int

type Error struct {
	pos string
	msg string
}

func (err *Error)Error() string  {
	return err.pos + err.msg
}

const (
	// Special tokens
	INT Token = iota
	CHAR
	FLOAT
	STRING
	VOID
	CLASS
	VIRTUAL
	EXTERN

	LPAREN
	LBRACK
	LBRACE
	COMMA
	PERIOD
	RPAREN
	RBRACK
	RBRACE
	SEMICOLON
	COLON
	PUBLIC
	PRIVATE
	PROTECTED
	USING
	MUL
	INVOKEFUNC
	INVOKEACTION
)

var tokens = [...]string{

	INT:     "int",
	FLOAT:   "float",
	CHAR:    "char",
	STRING:  "string",
	CLASS:   "class",
	VOID:    "void",
	VIRTUAL: "virtual",
	EXTERN: "extern",

	PUBLIC: "public",
	PRIVATE: "private",
	PROTECTED: "protected",
	USING: "using",
	INVOKEFUNC: "INVOKE_FUNC",
	INVOKEACTION: "INVOKE_ACTION",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",
	MUL: "*",
}
