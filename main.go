package main

import (
	"fmt"
	"log"

	"github.com/ayaki102/willNo/willno"
)

func main() {
	// API concept
	parser := willno.ParseFile("myLang.com"). // THIS IS BASE OF API
							AddKeyword("let").
							AddKeyword("fn").
							AddComment("//").
							AddComment("#").
							AddLiteral(willno.StringLiteral).
							AddLiteral(willno.NumberLiteral)

	parsed, err := parser.Parse()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(parsed)
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
