package willno

import (
	"go/token"
	"strings"
)

const (
	StringLiteral = "string"
	NumberLiteral = "number"
	BoolLiteral   = "boolean"
)

const (
	Var = "variable"
	Fn  = "function"
)

// LanguageBuilder: for chaining API
type LanguageBuilder struct {
	lang *language
}

// Core Value representation
type Value struct {
	Name        string // identifier name (variable, function, etc.)
	Type        string // token type: "variable", "function", "literal", "comment", etc.
	LiteralType string // "string", "number", "boolean", or custom literal type
	Value       any    // actual parsed value
}

// language: internal representation of a user-defined language
// end user won't see this
type language struct {
	name     string            // language identifier
	fileExts string            // optional: file extensions
	keywords map[string]string // logical type -> actual keyword ("variable" -> "let")
	comments []string          // comment patterns, e.g., "//", "#"
	literals []string          // supported literal types: "string", "number", "boolean", or custom
}

type Parsed struct {
	pf *parsedFile
}

// parsedFile: end-user access
type parsedFile struct {
	tokens map[string]map[string]Value // type -> name -> Value
}

func ParseFile(name string) *LanguageBuilder {
	filename := strings.Split(name, ".")
	return &LanguageBuilder{
		lang: &language{
			// only reasonable way to do this lol
			name:     filename[0],
			fileExts: filename[1],
			keywords: map[string]string{},
			comments: []string{},
			literals: []string{},
		},
	}
}

func (lb *LanguageBuilder) AddKeyword(literal, value string) *LanguageBuilder {
	lb.lang.keywords[literal] = value
	return lb
}

func (lb *LanguageBuilder) AddComment(comment string) *LanguageBuilder {
	lb.lang.comments = append(lb.lang.comments, comment)
	return lb
}

// this is used to declare types that are "supported" bu user's lang
func (lb *LanguageBuilder) AddLiteral(lit string) *LanguageBuilder {
	lb.lang.literals = append(lb.lang.literals, lit)
	return lb
}

func (lb *LanguageBuilder) Parse() *Parsed {
	// this will be from file
	// whole ass parser
	return &Parsed{
		pf: &parsedFile{
			tokens: map[string]map[string]Value{},
		},
	}
}

func (p *Parsed) Get(tokenType, name string) (any, bool) {
	return p.pf.Get(tokenType, name)
}

func (p *Parsed) GetString(tokenType, name string) (string, bool) {
	return p.pf.GetString(tokenType, name)
}

func (p *Parsed) GetNumber(tokenType, name string) (float64, bool) {
	return p.pf.GetNumber(tokenType, name)
}

func (p *Parsed) GetOr(tokenType, name string, def any) any {
	return p.pf.GetOr(tokenType, name, def)
}

func (p *Parsed) All(tokenType string) []Value {
	return p.pf.All(tokenType)
}

func (p *Parsed) Filter(tokenType string, predicate func(Value) bool) []Value {
	return p.pf.Filter(tokenType, predicate)
}

// private stuff that user won't see WHOLE IMPLEMENTATIONS
func (pf *parsedFile) Get(tokenType, name string) (any, bool) {
	return nil, false
}

func (pf *parsedFile) GetString(tokenType, name string) (string, bool) {
	val, ok := pf.Get(tokenType, name)
	if !ok {
		return "", false
	}
	s, ok := val.(string)
	return s, ok
}

func (pf *parsedFile) GetNumber(tokenType, name string) (float64, bool) {
	val, ok := pf.Get(tokenType, name)
	if !ok {
		return 0, false
	}
	n, ok := val.(float64)
	return n, ok
}

func (pf *parsedFile) GetOr(tokenType, name string, def any) any {
	if val, ok := pf.Get(tokenType, name); ok {
		return val
	}
	return def
}

func (pf *parsedFile) All(tokenType string) []Value {
	return nil
}

func (pf *parsedFile) Filter(tokenType string, predicate func(Value) bool) []Value {
	return nil
}
