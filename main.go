package main

func main() {
	// API concept
	parser := NewParser()

	parser.ReadFile("myLang"). // THIS IS BASE OF API
					AddKeyword("variable", "let").
					AddKeyword("function", "fn").
					AddComment("//").
					AddComment("#").
					AddLiteral("string").
					AddLiteral("number").
		//AddKeywordS add AddLiteralS for easier creation of multiple keywords in one call
		// optional: hooks for custom behavior
		OnKeywordMatch(func(tok Token) {
			fmt.Println("Matched keyword:", tok.Value)
		}).
		OnParseComplete(func(tokens []Token) {
			fmt.Println("Parsing finished! Total tokens:", len(tokens))
		})

}
