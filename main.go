package main

import (
	"github.com/ayaki102/willNo/willno"
)

func main() {
	// API concept
	parser := willno.ParseFile("myLang.com"). // THIS IS BASE OF API
							AddKeyword(willno.Var, "let").
							AddKeyword(willno.Fn, "fn").
							AddComment("//").
							AddComment("#").
							AddLiteral(willno.StringLiteral).
							AddLiteral(willno.NumberLiteral)

	parsed := parser.Parse() // for best should it return *ParsedFile struct?
	// or shoudl it also be wrraped so i can call parsed.Get() etc etc
	// parsed.Get("variable", "X")
}

//AddKeywordS add AddLiteralS for easier creation of multiple keywords in one call
// optional: hooks for custom behavior
// OnKeywordMatch(func(tok Token) {
// 	fmt.Println("Matched keyword:", tok.Value)
// }).
// OnParseComplete(func(tokens []Token) {
// 	fmt.Println("Parsing finished! Total tokens:", len(tokens))
// })
